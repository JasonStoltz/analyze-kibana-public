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
        "primaryText":   regexp.MustCompile(`colors.primaryText`),
        "accentText":    regexp.MustCompile(`colors.accentText`),
        "warningText":   regexp.MustCompile(`colors.warningText`),
        "dangerText":    regexp.MustCompile(`colors.dangerText`),
        "text":          regexp.MustCompile(`colors.text`),
        "title":         regexp.MustCompile(`colors.title`),
        "subduedText":   regexp.MustCompile(`colors.subduedText`),
        "disabledText":  regexp.MustCompile(`colors.disabledText`),
        "successText":   regexp.MustCompile(`colors.successText`),
    }

    utils.ProcessAndWriteResults(targetDirectory, "../kibana/.github/CODEOWNERS", "old_tokens_by_code_owner/results.csv", patterns)
}