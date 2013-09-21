package main

import (
	"flag"
	"fmt"
	"github.com/aybabtme/color/brush"
	"log"
	"os"
	"sort"
	"strings"
)

type statusCmd struct {
}

func (cmd *statusCmd) Name() string {
	return "status"
}

func (cmd *statusCmd) DefineFlags(fs *flag.FlagSet) {
}

func (cmd *statusCmd) Run(fs *flag.FlagSet) {
	if fs.NArg() == 0 {
		for _, repo := range SortedRealRepos {
			printGreen("%s\n", repo.Path())
			printStatus(repo)
		}
	} else {
		for _, key := range fs.Args() {
			repo := RealRepos[PaniniTag(key)]
			printGreen("%s\n", repo.Path())
			printStatus(repo)
		}
	}
}

func printStatus(r *RealRepo) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	os.Chdir(r.path)
	defer os.Chdir(cwd)
	for _, line := range strings.Split(r.GitExec("status", "--short"), "\n") {
		if line == "" {
			continue
		}
		printRed("\t%s\n", line)
	}

	remotes := []string{}
	for key, _ := range r.Remotes() {
		remotes = append(remotes, key)
	}
	sort.Strings(remotes)
	remoteBranches := r.RemoteBranches()
	localBranches := r.LocalBranches()
	sort.Strings(localBranches)
	for _, branch := range localBranches {
		fmt.Printf("\t%s", branch)
		for _, remote := range remotes {
			if _, ok := remoteBranches[remote]; !ok {
				fmt.Printf(brush.Red(" [%s nobranch]").String(), remote)
			} else if _, ok := remoteBranches[remote][branch]; !ok {
				fmt.Printf(brush.Red(" [%s nobranch]").String(), remote)
			} else {
				diff := r.Difference(remote+"/"+branch, branch)
				if !diff.Diverged && diff.Difference == 0 {
					fmt.Printf(" [%s %s]", remote, diff)
				} else {
					fmt.Printf(brush.Red(" [%s %s]").String(), remote, diff)
				}
			}
		}
		fmt.Printf("\n")
	}
}
