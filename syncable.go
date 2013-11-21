package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
)

type syncableCmd struct {
}

func (cmd *syncableCmd) Name() string {
	return "syncable"
}

func (cmd *syncableCmd) DefineFlags(fs *flag.FlagSet) {
}

func (cmd *syncableCmd) Run(fs *flag.FlagSet) {
	output, err := exec.Command("git", "rev-parse", "--show-toplevel").Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "It seems the cwd is not a git repository")
		os.Exit(1)
	}
	repoPath := strings.TrimSpace(string(output))
	paniniTag := PaniniTag("panini:" + path.Base(repoPath))

	realRepo := RealRepo{repoPath, paniniTag}
	if paniniRepo := realRepo.PaniniRepo(); paniniRepo != nil {
		if paniniRepo.PaniniRepoPath() != "" {
			fmt.Fprintf(os.Stderr, "It seems %q is already syncable\n", realRepo.tag)
			os.Exit(1)
		}
	}
	if err := CreatePaniniRepo(realRepo); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Clone to SyncBase done\n")
}
