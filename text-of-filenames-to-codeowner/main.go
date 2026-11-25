package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/hmarr/codeowners"
	"github.com/samber/lo"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	codeownersFile, err := os.Open("../kibana/.github/CODEOWNERS")
	check(err)

	codeownersRuleset, err := codeowners.ParseFile(codeownersFile)
	check(err)

	file, err := os.Open("./text-of-filenames-to-codeowner/tmp.txt")
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fileName := scanner.Text()
		codeownerRule, err := codeownersRuleset.Match(fileName)
		check(err)

		if codeownerRule != nil {
			owners := lo.Map(codeownerRule.Owners, func(owner codeowners.Owner, index int) string {
				return owner.Value
			})

			fmt.Println(strings.Join(owners, ","))
		} else {
			fmt.Println("null")
		}
	}
}
