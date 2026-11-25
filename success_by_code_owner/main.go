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
        "color=\"success\"":                regexp.MustCompile(`color="success"`),
        "$euiColorSuccess":                 regexp.MustCompile(`\$euiColorSuccess`),
        "$euiColorSuccessText":             regexp.MustCompile(`\$euiColorSuccessText`),
        "euiTheme.colors.success":          regexp.MustCompile(`euiTheme\.colors\.success`),
        "euiTheme.colors.successText":      regexp.MustCompile(`euiTheme\.colors\.successText`),
        "theme.eui.euiColorSuccess":        regexp.MustCompile(`theme\.eui\.euiColorSuccess`),
        "theme.eui.euiColorSuccessText":    regexp.MustCompile(`theme\.eui\.euiColorSuccessText`),
    }

    utils.ProcessAndWriteResults(targetDirectory, "../kibana/.github/CODEOWNERS", "success_by_code_owner/results.csv", patterns)
}