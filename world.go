package main

import (
	"flag"
	"fmt"
)

type worldCmd struct {
	verbose bool
	local   bool
}

func (cmd *worldCmd) Name() string {
	return "world"
}

func (cmd *worldCmd) DefineFlags(fs *flag.FlagSet) {
	fs.BoolVar(&cmd.verbose, "verbose", false, "more output")
	fs.BoolVar(&cmd.local, "local", false, "show only local")
}

func (cmd *worldCmd) Run(fs *flag.FlagSet) {
	for _, repo := range SortedRealRepos {
		if cmd.local {
			if _, ok := RealRepos[repo.Tag()]; !ok {
				continue
			}
		}
		fmt.Println(repo.Tag())
		if !cmd.verbose {
			continue
		}
		for name, remote := range repo.Remotes() {
			fmt.Printf("\t%s %s (fetch)\n", name, remote.FetchURL)
			for _, url := range remote.PushURLs {
				fmt.Printf("\t%s %s (push)\n", name, url)
			}
		}
	}
}
