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
        "euiThemeVars": regexp.MustCompile(`euiThemeVars.\w+`),
    }

    utils.ProcessAndWriteResults(targetDirectory, "../kibana/.github/CODEOWNERS", "eui_theme_vars_by_code_owner/results.csv", patterns)
}