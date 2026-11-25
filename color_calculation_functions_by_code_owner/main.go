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
        "shade()":          regexp.MustCompile(`shade\(`),
        "tint()":           regexp.MustCompile(`tint\(`),
        "shadeOrTint()":    regexp.MustCompile(`shadeOrTint\(`),
        "tintOrShade()":    regexp.MustCompile(`tintOrShade\(`),
        "transparentize()": regexp.MustCompile(`transparentize\(`),
    }

    utils.ProcessAndWriteResults(targetDirectory, "../kibana/.github/CODEOWNERS", "color_calculation_functions_by_code_owner/results.csv", patterns)
}