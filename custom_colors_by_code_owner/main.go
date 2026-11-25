package main

import (
	"os"
	"regexp"

	"analyze-kibana/utils"
)

func main() {
	targetDirectory := "../kibana"
	if len(os.Args) > 1 {
		targetDirectory = os.Args[1]
	}

	patterns := map[string]*regexp.Regexp{
		"HEX":  regexp.MustCompile(`#[0-9a-fA-F]{3,6}`),
		"RGB":  regexp.MustCompile(`rgb\(\s*\d{1,3}\s*,\s*\d{1,3}\s*,\s*\d{1,3}\s*\)`),
		"RGBA": regexp.MustCompile(`rgba\(\s*\d{1,3}\s*,\s*\d{1,3}\s*,\s*\d{1,3}\s*,\s*(0|1|0?\.\d+)\s*\)`),
	}

	utils.ProcessAndWriteResults(targetDirectory, "../kibana/.github/CODEOWNERS", "custom_colors_by_code_owner/results.csv", patterns)
}
