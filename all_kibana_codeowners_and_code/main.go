package main

import (
	"fmt"
	"log"
	"os"

	"github.com/hmarr/codeowners"
	"github.com/samber/lo"
)

func main() {
	file, err := os.Open("../kibana/.github/CODEOWNERS")
	if err != nil {
		log.Fatal(err)
	}

	ruleset, err := codeowners.ParseFile(file)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Owner\tFiles")
	ruleStrings := lo.FlatMap(ruleset, func(rule codeowners.Rule, index int) []string {
		return lo.Map(rule.Owners, func(owner codeowners.Owner, index int) string {
			return owner.Value + "\t" + rule.RawPattern()
		})
	})

	lo.ForEach(ruleStrings, func(ruleString string, index int) {
		fmt.Println(ruleString)
	})
}
