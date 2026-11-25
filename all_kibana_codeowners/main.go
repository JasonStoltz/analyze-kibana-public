package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

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

	list := getListOfOwners(ruleset)

	fmt.Println("Owner\tLink")
	lo.ForEach(list, func(owner string, index int) {
		fmt.Println(owner + "\thttps://github.com/orgs/elastic/teams/" + strings.Replace(owner, "elastic/", "", 1))
	})
}

func getListOfOwners(ruleset codeowners.Ruleset) []string {

	allOwners := lo.FlatMap(ruleset, func(rule codeowners.Rule, index int) []string {
		return lo.Map(rule.Owners, func(owner codeowners.Owner, index int) string {
			return owner.Value
		})
	})

	return sort.StringSlice(lo.Uniq(allOwners))
}
