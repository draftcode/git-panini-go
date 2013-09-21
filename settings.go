package main

import (
	"io/ioutil"
	"launchpad.net/goyaml"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var (
	SyncBase        string
	RealRepos       map[PaniniTag]*RealRepo = make(map[PaniniTag]*RealRepo)
	SortedRealRepos []*RealRepo
	PaniniRepos     map[PaniniTag]*PaniniRepo = make(map[PaniniTag]*PaniniRepo)
	IgnorePaths     map[string]bool           = make(map[string]bool)
)

type localSettings struct {
	SyncBase     string
	Repositories []*RealRepo
	Ignore       []string
}

type paniniSettings struct {
	Repositories []*PaniniRepo
}

func listPaniniTagsFromSyncBase() []PaniniTag {
	file, err := os.Open(SyncBase)
	if err != nil {
		log.Fatal(err)
	}
	names, err := file.Readdirnames(-1)
	if err != nil {
		log.Fatal(err)
	}

	ret := []PaniniTag{}
	for _, name := range names {
		if strings.HasSuffix(name, ".git") {
			ret = append(ret, PaniniTag("panini:"+name[:len(name)-4]))
		}
	}
	return ret
}

func ReadSettings() {
	{
		b, err := ioutil.ReadFile(os.ExpandEnv("$HOME/.git-panini"))
		if err != nil {
			log.Fatal(err)
		}
		var settings localSettings
		if err := goyaml.Unmarshal(b, &settings); err != nil {
			log.Fatal(err)
		}

		SyncBase = fixPath(settings.SyncBase)
		for _, repo := range settings.Repositories {
			RealRepos[repo.Tag()] = repo
		}
		for _, path := range settings.Ignore {
			IgnorePaths[fixPath(path)] = true
		}

		tags := PaniniTags{}
		for tag, _ := range RealRepos {
			tags = append(tags, tag)
		}
		sort.Sort(&tags)
		for _, tag := range tags {
			SortedRealRepos = append(SortedRealRepos, RealRepos[tag])
		}
	}

	if SyncBase != "" {
		b, err := ioutil.ReadFile(filepath.Join(SyncBase, "panini-info"))
		if err != nil {
			log.Fatal(err)
		}
		var settings paniniSettings
		if err := goyaml.Unmarshal(b, &settings); err != nil {
			log.Fatal(err)
		}

		for _, repo := range settings.Repositories {
			PaniniRepos[repo.Tag()] = repo
		}

		for _, tagName := range listPaniniTagsFromSyncBase() {
			if _, ok := PaniniRepos[tagName]; !ok {
				PaniniRepos[tagName] = NewPaniniRepo(tagName)
			}
		}
	}
}
