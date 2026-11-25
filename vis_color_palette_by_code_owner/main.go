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
        "euiColorVis0": regexp.MustCompile(`euiColorVis0`),
        "euiColorVis1": regexp.MustCompile(`euiColorVis1`),
        "euiColorVis2": regexp.MustCompile(`euiColorVis2`),
        "euiColorVis3": regexp.MustCompile(`euiColorVis3`),
        "euiColorVis4": regexp.MustCompile(`euiColorVis4`),
        "euiColorVis5": regexp.MustCompile(`euiColorVis5`),
        "euiColorVis6": regexp.MustCompile(`euiColorVis6`),
        "euiColorVis7": regexp.MustCompile(`euiColorVis7`),
        "euiColorVis8": regexp.MustCompile(`euiColorVis8`),
        "euiColorVis9": regexp.MustCompile(`euiColorVis9`),
    }

    utils.ProcessAndWriteResults(targetDirectory, "../kibana/.github/CODEOWNERS", "vis_color_palette_by_code_owner/results.csv", patterns)
}