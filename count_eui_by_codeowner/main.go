package main

import (
	"analyze-kibana/utils"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/gobwas/glob"
	"github.com/hmarr/codeowners"
	"github.com/samber/lo"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	file, err := os.Open("../kibana/.github/CODEOWNERS")
	check(err)

	ruleset, err := codeowners.ParseFile(file)
	check(err)

	var numOfFiles int
	allFiles := make(map[string]int)
	tsxFiles := make(map[string]int)
	euiFiles := make(map[string]int)

	targetDirectory := "../kibana"
	if len(os.Args) > 1 {
		targetDirectory = os.Args[1]
	}
	textToFind := "<Eui"
	if len(os.Args) > 2 {
		textToFind = os.Args[2]
	}

	filepath.WalkDir(targetDirectory, func(s string, d fs.DirEntry, err error) error {
		check(err)

		includeList := glob.MustCompile("{*.ts,*.js,*.tsx,*.scss,*.yaml,*.html}")
		excludeList := glob.MustCompile("{*examples/*,*.test*,*stories*,*.asciidoc,*.ndjson,*.mdx,*test*,*.md,*.svg,*.lock}")

		if !d.IsDir() {
			numOfFiles++
			fmt.Fprintf(os.Stderr, "\rOn file %d", numOfFiles)

			rule, err := ruleset.Match(strings.Replace(s, "../kibana/", "", 1))
			check(err)

			if rule != nil && includeList.Match(s) && !excludeList.Match(s) {
				euiCount, err := utils.CountGrepMatches("-oP", textToFind, s)
				check(err)
				for _, owner := range rule.Owners {
					allFiles[owner.Value]++
					if strings.HasSuffix(s, ".tsx") {
						tsxFiles[owner.Value]++
					}
					euiFiles[owner.Value] += euiCount
				}
			}
		}

		return nil
	})

	fmt.Fprintln(os.Stderr)
	keys := lo.Keys(allFiles)
	fmt.Println("Codeowner" + "\t" + "All Owned Files" + "\t" + "All Owned .tsx Files" + "\t" + "All EUI usages")
	lo.ForEach(keys, func(key string, index int) {
		fmt.Println(key + "\t" + fmt.Sprint(allFiles[key]) + "\t" + fmt.Sprint(tsxFiles[key]) + "\t" + fmt.Sprint(euiFiles[key]))
	})
}
