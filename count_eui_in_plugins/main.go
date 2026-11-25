package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"analyze-kibana/utils"
	"github.com/hmarr/codeowners"
)

const kibana = "../kibana"

func pluginDirs() ([]string, error) {
	var plugins []string
	legacy := filepath.Join(kibana, "x-pack/plugins")
	entries, err := os.ReadDir(legacy)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}
	for _, entry := range entries {
		if entry.IsDir() {
			plugins = append(plugins, filepath.Join(legacy, entry.Name()))
		}
	}

	for _, pattern := range []string{"x-pack/platform/plugins", "x-pack/solutions/*/plugins"} {
		roots, err := filepath.Glob(filepath.Join(kibana, pattern))
		if err != nil {
			return nil, err
		}
		for _, root := range roots {
			err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
				if err != nil || !entry.IsDir() {
					return err
				}
				_, err = os.Stat(filepath.Join(path, "kibana.jsonc"))
				if os.IsNotExist(err) {
					return nil
				}
				if err != nil {
					return err
				}
				plugins = append(plugins, path)
				return filepath.SkipDir
			})
			if err != nil {
				return nil, err
			}
		}
	}
	if len(plugins) == 0 {
		return nil, fmt.Errorf("no X-Pack plugins found in %s", kibana)
	}
	return plugins, nil
}

func main() {
	file, err := os.Open(filepath.Join(kibana, ".github/CODEOWNERS"))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	ruleset, err := codeowners.ParseFile(file)
	if err != nil {
		log.Fatal(err)
	}

	plugins, err := pluginDirs()
	if err != nil {
		log.Fatal(err)
	}
	pattern := "Eui"
	if len(os.Args) > 1 {
		pattern = os.Args[1]
	}

	fmt.Println("Plugin\tEUI usages\tCodeowners")
	for _, path := range plugins {
		count, err := utils.CountGrepMatches("-Iroh", pattern, path)
		if err != nil {
			log.Fatal(err)
		}
		rel, err := filepath.Rel(kibana, path)
		if err != nil {
			log.Fatal(err)
		}
		rule, err := ruleset.Match(filepath.ToSlash(rel))
		if err != nil {
			log.Fatal(err)
		}
		var owners []string
		if rule != nil {
			for _, owner := range rule.Owners {
				owners = append(owners, owner.Value)
			}
		}
		fmt.Printf("%s\t%d\t%s\n", filepath.Base(path), count, strings.Join(owners, ","))
	}
}
