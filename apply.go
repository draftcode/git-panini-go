package main

import (
	"flag"
	"fmt"
	"os"
)

type applyCmd struct {
	verbose bool
	force   bool
}

func (cmd *applyCmd) Name() string {
	return "apply"
}

func (cmd *applyCmd) DefineFlags(fs *flag.FlagSet) {
	fs.BoolVar(&cmd.verbose, "verbose", false, "more output")
	fs.BoolVar(&cmd.force, "force", false, "forcibly set remotes")
}

func (cmd *applyCmd) Run(fs *flag.FlagSet) {
	for _, repo := range SortedRealRepos {
		apply(repo, true, cmd.verbose, cmd.force)
	}
}

type noopCmd struct {
	verbose bool
	force   bool
}

func (cmd *noopCmd) Name() string {
	return "noop"
}

func (cmd *noopCmd) DefineFlags(fs *flag.FlagSet) {
	fs.BoolVar(&cmd.verbose, "verbose", false, "more output")
	fs.BoolVar(&cmd.force, "force", false, "forcibly set remotes")
}

func (cmd *noopCmd) Run(fs *flag.FlagSet) {
	for _, repo := range SortedRealRepos {
		apply(repo, false, cmd.verbose, cmd.force)
	}
}

func apply(repo *RealRepo, effective, verbose, force bool) {
	printGreen("Checking %s...\n", repo.Path())
	paniniRepo := repo.PaniniRepo()
	if paniniRepo == nil {
		printRed("\tIt has an invalid panini repository: %s\n", repo.Tag())
		return
	}
	if _, err := os.Stat(repo.Path()); os.IsNotExist(err) {
		printRed("\tIt is not checked out.")
		if effective {
			checkout(repo)
		} else {
			// We cannot do anything if we don't have a repository.
			return
		}
	}
	if repo.GitExec("rev-parse", "--is-inside-work-tree") != "true\n" {
		printRed("\tIt is not a work tree")
		return
	}
	compareRemotes(repo, paniniRepo, effective, verbose, force)
}

func checkout(repo *RealRepo) {
}

func compareRemotes(repo *RealRepo, paniniRepo *PaniniRepo,
	effective, verbose, force bool) {
	current := repo.Remotes()
	setting := paniniRepo.Remotes()

	for name, remote := range setting {
		if currentRemote, ok := current[name]; ok {
			if remote.FetchURL != currentRemote.FetchURL {
				printRed("\tDifference current: %s %s (fetch)\n",
					name, currentRemote.FetchURL)
				printRed("\t           setting: %s %s (fetch)\n",
					name, remote.FetchURL)
				if force {
					printRed("\tForcibly set url\n")
					if effective {
						repo.GitExec("remote", "set-url", name, remote.FetchURL)
					}
				}
			} else {
				if verbose {
					fmt.Printf("\tSame: %s %s (fetch)\n", name, remote.FetchURL)
				}
			}
			m := make(map[string]bool)
			for _, url := range currentRemote.PushURLs {
				m[url] = false
			}
			for _, url := range remote.PushURLs {
				if _, ok := m[url]; ok {
					m[url] = true
					if verbose {
						fmt.Printf("\tSame: %s %s (push)\n", name, url)
					}
				} else {
					printRed("\tSet: %s %s (push)\n", name, url)
					if effective {
						repo.GitExec("remote", "set-url", "--add", "--push",
							name, url)
					}
				}
			}
			for url, checked := range m {
				if !checked {
					printRed("\tNot in setting. Ignore: %s %s (push)\n",
						name, url)
				}
			}
		} else {
			printRed("\tSet: %s %s (fetch)\n", name, remote.FetchURL)
			if effective {
				repo.GitExec("remote", "add", name, remote.FetchURL)
			}
			for _, url := range remote.PushURLs {
				printRed("\tSet: %s %s (push)\n", name, url)
				if effective {
					repo.GitExec("remote", "set-url", "--add", "--push",
						name, url)
				}
			}
		}
	}

	for name, _ := range current {
		if _, ok := setting[name]; !ok {
			printRed("\tNot in setting. Ignore: %s\n", name)
		}
	}
}
