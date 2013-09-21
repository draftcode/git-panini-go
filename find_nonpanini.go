package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type findNonPaniniCmd struct {
}

func (cmd *findNonPaniniCmd) Name() string {
	return "find-nonpanini"
}

func (cmd *findNonPaniniCmd) DefineFlags(fs *flag.FlagSet) {
}

func (cmd *findNonPaniniCmd) Run(fs *flag.FlagSet) {
	m := make(map[string]bool)
	for _, repo := range RealRepos {
		m[repo.Path()] = true
	}

	if fs.NArg() == 0 {
		cwd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		for _, path := range findRepos(cwd) {
			if _, ok := m[path]; !ok {
				fmt.Println(path)
			}
		}
	} else {
		for _, arg := range fs.Args() {
			arg, err := filepath.Abs(arg)
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, path := range findRepos(arg) {
				if _, ok := m[path]; !ok {
					fmt.Println(path)
				}
			}
		}
	}
}

func findRepos(dir string) []string {
	results := []string{}
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		if _, ok := IgnorePaths[path]; ok {
			return filepath.SkipDir
		}

		fi, err := os.Stat(filepath.Join(path, ".git"))
		if os.IsNotExist(err) || fi == nil {
			return nil
		}
		if fi.IsDir() {
			results = append(results, path)
			return filepath.SkipDir
		} else {
			return nil
		}
	})
	if err != nil {
		log.Fatal(err)
	}
	return results
}
