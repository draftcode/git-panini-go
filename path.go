package main

import (
	"flag"
	"fmt"
	"os"
)

type pathCmd struct {
}

func (cmd *pathCmd) Name() string {
	return "path"
}

func (cmd *pathCmd) DefineFlags(fs *flag.FlagSet) {
}

func (cmd *pathCmd) Run(fs *flag.FlagSet) {
	if fs.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "path command takes only one argument")
		os.Exit(1)
	}
	tag := PaniniTag(fs.Arg(0))
	repo, ok := RealRepos[tag]
	if !ok {
		fmt.Fprintf(os.Stderr, "Repository %s is not found\n", tag)
		os.Exit(1)
	}
	fmt.Println(repo.Path())
}
