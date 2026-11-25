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
	themeVarFiles := make(map[string]int)
	customColorFiles := make(map[string]int)

	targetDirectory := "../kibana"
	if len(os.Args) > 1 {
		targetDirectory = os.Args[1]
	}

	euiUsageGrep := "<Eui"
	euiThemeVarsUsageGrep := "euiThemeVars"
	// This is a regex suggested by chatgpt to catch 3 digit and 6 digits hex codes, rgb, hsl.
	// Temporary note, to test directly:
	//    ggrep -roP "#([0-9a-fA-F]{6}|[0-9a-fA-F]{3})|rgb(a)?\([0-9]+,\s*[0-9]+,\s*[0-9]+(,\s*[0-9.]+)?\)|hsl(a)?\([0-9]+,\s*[0-9]+%,\s*[0-9]+%(,\s*[0-9.]+)?\)\b" ../kibana/x-pack/plugins/enterprise_search --exclude="*.test*" --exclude="*.asciidoc" --exclude="*.ndjson" --exclude="*.mdx" --exclude="*.md" --exclude="*.svg" --exclude="*.lock"  --exclude="*.chunk*" --exclude="*.map*" --exclude="*.png*" --exclude="*.mock*" --exclude-dir="*examples/*" --exclude-dir="*stories*" --exclude-dir="*test*"
	customerColorGrep := `#([0-9a-fA-F]{6}|[0-9a-fA-F]{3})|rgb(a)?\([0-9]+,\s*[0-9]+,\s*[0-9]+(,\s*[0-9.]+)?\)|hsl(a)?\([0-9]+,\s*[0-9]+%,\s*[0-9]+%(,\s*[0-9.]+)?\)\b`

	filepath.WalkDir(targetDirectory, func(s string, d fs.DirEntry, err error) error {
		check(err)

		includeList := glob.MustCompile("{*.ts,*.js,*.tsx,*.scss,*.yaml,*.html}")
		excludeList := glob.MustCompile("{*examples/*,*.test*,*stories*,*.asciidoc,*.ndjson,*.mdx,*test*,*.md,*.svg,*.lock,*.chunk*,*.map*,*.png}")

		if !d.IsDir() {
			numOfFiles++
			fmt.Fprintf(os.Stderr, "\rOn file %d", numOfFiles)

			rule, err := ruleset.Match(strings.Replace(s, "../kibana/", "", 1))
			check(err)

			if rule != nil && includeList.Match(s) && !excludeList.Match(s) {
				euiCount, err := utils.CountGrepMatches("-oP", euiUsageGrep, s)
				check(err)
				themeCount, err := utils.CountGrepMatches("-oP", euiThemeVarsUsageGrep, s)
				check(err)
				colorCount, err := utils.CountGrepMatches("-oP", customerColorGrep, s)
				check(err)
				for _, owner := range rule.Owners {
					allFiles[owner.Value]++
					if strings.HasSuffix(s, ".tsx") {
						tsxFiles[owner.Value]++
					}
					euiFiles[owner.Value] += euiCount
					themeVarFiles[owner.Value] += themeCount
					customColorFiles[owner.Value] += colorCount
				}
			}
		}

		return nil
	})

	fmt.Fprintln(os.Stderr)
	keys := lo.Keys(allFiles)
	fmt.Println("Codeowner" + "\t" + "All Owned Files" + "\t" + "All Owned .tsx Files" + "\t" + "All EUI usages" + "\t" + "euiThemeVars Usages" + "\t" + "custom colors Usages")
	lo.ForEach(keys, func(key string, index int) {
		fmt.Println(key + "\t" + fmt.Sprint(allFiles[key]) + "\t" + fmt.Sprint(tsxFiles[key]) + "\t" + fmt.Sprint(euiFiles[key]) + "\t" + fmt.Sprint(themeVarFiles[key]) + "\t" + fmt.Sprint(customColorFiles[key]))
	})
}
