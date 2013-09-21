package main

import (
	"flag"
	"fmt"
	"github.com/aybabtme/color/brush"
)

type fetchCmd struct {
}

func (cmd *fetchCmd) Name() string {
	return "fetch"
}

func (cmd *fetchCmd) DefineFlags(fs *flag.FlagSet) {
}

func (cmd *fetchCmd) Run(fs *flag.FlagSet) {
	fmt.Printf(brush.Green("Execute git fetch in all repositories...\n").String())
	for _, repo := range SortedRealRepos {
		fmt.Println(repo.Path())
		for name, _ := range repo.Remotes() {
			fmt.Println("\t", name)
			repo.Fetch(name)
		}
	}
	fmt.Printf(brush.Green("Complete!\n").String())
}
