package main

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var (
	CannotCreateRepo = errors.New("Cannot create a repository")
)

type (
	PaniniTag  string
	PaniniTags []PaniniTag
)

func (t PaniniTag) Name() string {
	return strings.Replace(string(t), "panini:", "", 1)
}

func (t *PaniniTags) Len() int {
	return len(*t)
}

func (t *PaniniTags) Less(i, j int) bool {
	return (*t)[i] < (*t)[j]
}

func (t *PaniniTags) Swap(i, j int) {
	temp := (*t)[i]
	(*t)[i] = (*t)[j]
	(*t)[j] = temp
}

type PaniniRepo struct {
	tag     PaniniTag
	remotes map[string]Remote
}

func (r *PaniniRepo) SetYAML(tag string, value interface{}) bool {
	if tag != "!!map" {
		return false
	}
	m := value.(map[interface{}]interface{})
	tagName, ok := m["name"]
	if !ok {
		return false
	}
	r.tag = PaniniTag(tagName.(string))
	r.remotes = make(map[string]Remote)
	if value, ok := m["remotes"]; ok {
		remotes := value.(map[interface{}]interface{})
		for name, v := range remotes {
			r.remotes[name.(string)] = r.parseRemote(v)
		}
	}
	return true
}

func (r *PaniniRepo) parseRemote(value interface{}) Remote {
	switch v := value.(type) {
	case string:
		return Remote{FetchURL: v, PushURLs: []string{v}}
	case map[interface{}]interface{}:
		ret := Remote{}
		if url, ok := v["url"]; ok {
			ret.FetchURL = url.(string)
			ret.PushURLs = append(ret.PushURLs, url.(string))
		}
		if url, ok := v["fetch"]; ok {
			ret.FetchURL = url.(string)
		}
		if urls, ok := v["push"]; ok {
			if url, ok := urls.(string); ok {
				ret.PushURLs = append(ret.PushURLs, url)
			} else if urls, ok := urls.([]interface{}); ok {
				for _, url := range urls {
					ret.PushURLs = append(ret.PushURLs,
						url.(string))
				}
			}
		}
		return ret
	default:
		panic(value)
	}
}

func NewPaniniRepo(tag PaniniTag) *PaniniRepo {
	r := &PaniniRepo{
		tag:     tag,
		remotes: make(map[string]Remote),
	}
	return r
}

func (r *PaniniRepo) PaniniRepoPath() string {
	if SyncBase == "" {
		return ""
	}
	path := filepath.Join(SyncBase, r.Name()+".git")
	if fi, err := os.Stat(path); os.IsNotExist(err) || !fi.IsDir() {
		return ""
	}
	output, _ := exec.Command("git", "--git-dir", path, "rev-parse", "--is-bare-repository").Output()
	if string(output) != "true\n" {
		return ""
	}
	return path
}

// CreatePaniniRepo creates a Panini repository from a real repository.
func CreatePaniniRepo(realRepo RealRepo) error {
	if SyncBase == "" {
		return CannotCreateRepo
	}
	path := filepath.Join(SyncBase, realRepo.tag.Name()+".git")
	if err := exec.Command("git", "clone", "--mirror", realRepo.path, path).Run(); err != nil {
		return CannotCreateRepo
	}
	if err := exec.Command("git", "--git-dir", path, "remote", "rm", "origin").Run(); err != nil {
		return CannotCreateRepo
	}
	return nil
}

// Clone executes git clone to the specified filepath.
// func (r *PaniniRepo) Clone(path string) error {
// 	if _, err := os.Stat(path); os.IsExist(err) {
// 		return err
// 	}
// 	if err := exec.Command(
// 		"git", "clone", "--origin", "panini", "--", r.Path(), path,
// 	).Run(); err != nil {
// 		return err
// 	}
// 	return nil
// }

func (r *PaniniRepo) Name() string {
	return r.tag.Name()
}

func (r *PaniniRepo) Tag() PaniniTag {
	return r.tag
}

func (r *PaniniRepo) Remotes() map[string]Remote {
	remotes := make(map[string]Remote)
	for k, v := range r.remotes {
		remotes[k] = v
	}
	if path := r.PaniniRepoPath(); path != "" {
		remotes["panini"] = Remote{
			FetchURL: path,
			PushURLs: []string{path},
		}
	}
	return remotes
}
