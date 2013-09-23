package main

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	spacesPattern = regexp.MustCompile(`\s+`)
)

type Remote struct {
	FetchURL string
	PushURLs []string
}

type Difference struct {
	Diverged   bool
	Difference int
}

func (d Difference) String() string {
	if d.Diverged {
		return "diverged"
	} else {
		return fmt.Sprintf("%+d", d.Difference)
	}
}

// RealRepo is a repository that is not under the syncbase directory.
type RealRepo struct {
	path string
	tag  PaniniTag
}

func (r *RealRepo) SetYAML(tag string, value interface{}) bool {
	switch tag {
	case "!!str":
		r.path = fixPath(value.(string))
		r.tag = PaniniTag("panini:" + filepath.Base(r.path))
	case "!!map":
		m := value.(map[interface{}]interface{})
		if v, ok := m["path"]; ok {
			r.path = fixPath(v.(string))
		}
		if v, ok := m["panini"]; ok {
			r.tag = PaniniTag(v.(string))
		} else {
			r.tag = PaniniTag("panini:" + filepath.Base(r.path))
		}
	default:
		return false
	}
	return true
}

func (r *RealRepo) PaniniRepo() *PaniniRepo {
	if repo, ok := PaniniRepos[r.tag]; ok {
		return repo
	} else {
		return nil
	}
}

func (r *RealRepo) Path() string {
	return r.path
}

func (r *RealRepo) Tag() PaniniTag {
	return r.tag
}

func (r *RealRepo) GitExec(commands ...string) string {
	output, err := r.GitExecErr(commands...)
	if err != nil {
		log.Fatal(err)
	}
	return output
}

func (r *RealRepo) GitExecErr(commands ...string) (string, error) {
	args := []string{"--git-dir", filepath.Join(r.path, ".git")}
	args = append(args, commands...)
	output, err := exec.Command("git", args...).CombinedOutput()
	return string(output), err
}

func (r *RealRepo) Remotes() map[string]Remote {
	fetches := make(map[string]string)
	pushes := make(map[string][]string)
	for _, line := range strings.Split(r.GitExec("remote", "-vv"), "\n") {
		parts := spacesPattern.Split(line, -1)
		if len(parts) != 3 {
			continue
		}
		if parts[2] == "(fetch)" {
			fetches[parts[0]] = parts[1]
		} else if parts[2] == "(push)" {
			pushes[parts[0]] = append(pushes[parts[0]], parts[1])
		}
	}

	remotes := make(map[string]Remote)
	for k, fetchURL := range fetches {
		remotes[k] = Remote{FetchURL: fetchURL, PushURLs: pushes[k]}
	}
	return remotes
}

func (r *RealRepo) LocalBranches() (branches []string) {
	for _, line := range strings.Split(r.GitExec("branch"), "\n") {
		if len(line) < 2 {
			continue
		}
		line = line[2:]
		branches = append(branches, line)
	}
	return
}

func (r *RealRepo) RemoteBranches() map[string]map[string]bool {
	// map[remote][branch] = true
	branches := make(map[string]map[string]bool)
	for _, line := range strings.Split(r.GitExec("branch", "--remotes"), "\n") {
		if len(line) < 2 {
			continue
		}
		line = line[2:]
		parts := strings.Split(line, "/")
		remote := parts[0]
		branch := parts[1]

		if _, ok := branches[remote]; !ok {
			branches[remote] = make(map[string]bool)
		}
		branches[remote][branch] = true
	}
	return branches
}

func (r *RealRepo) Fetch(name string) {
	r.GitExec("fetch", name)
}

type Rel int

const (
	Behind1 Rel = iota
	Behind2
	Diverged
)

func (r *RealRepo) Relation(branch1, branch2 string) Rel {
	output, err := r.GitExecErr("merge-base", branch1, branch2)
	if err != nil {
		return Diverged
	}

	commit1 := r.GitExec("rev-parse", branch1)
	commit2 := r.GitExec("rev-parse", branch2)
	switch output {
	case commit1:
		return Behind1
	case commit2:
		return Behind2
	default:
		return Diverged
	}
}

// Difference calculates difference between local branch and remote branch.
func (r *RealRepo) Difference(branch1, branch2 string) Difference {
	rel := r.Relation(branch1, branch2)
	if rel == Diverged {
		return Difference{Diverged: true}
	}
	diff := 0
	if rel == Behind1 {
		diff = len(strings.Split(r.GitExec("rev-list", branch1+".."+branch2), "\n")) - 1
	} else {
		diff = -len(strings.Split(r.GitExec("rev-list", branch2+".."+branch1), "\n")) + 1
	}
	return Difference{Diverged: false, Difference: diff}
}
