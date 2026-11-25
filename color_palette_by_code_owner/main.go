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
        "theme.eui.*":                     regexp.MustCompile(`theme\.eui\..*`),
        "theme.eui.euiColorVis*":          regexp.MustCompile(`theme\.eui\.euiColorVis`),
        "euiPaletteColorBlind":            regexp.MustCompile(`euiPaletteColorBlind`),
        "euiPaletteColorBlindBehindText":  regexp.MustCompile(`euiPaletteColorBlindBehindText`),
        "euiPaletteForStatus":             regexp.MustCompile(`euiPaletteForStatus`),
        "euiPaletteForTemperature":        regexp.MustCompile(`euiPaletteForTemperature`),
        "euiPaletteComplementary":         regexp.MustCompile(`euiPaletteComplementary`),
        "euiPaletteRed":                   regexp.MustCompile(`euiPaletteRed`),
        "euiPaletteGreen":                 regexp.MustCompile(`euiPaletteGreen`),
        "euiPaletteCool":                  regexp.MustCompile(`euiPaletteCool`),
        "euiPaletteWarm":                  regexp.MustCompile(`euiPaletteWarm`),
        "euiPaletteGray":                  regexp.MustCompile(`euiPaletteGray`),
        "colorPalette":                    regexp.MustCompile(`colorPalette`),
    }

    utils.ProcessAndWriteResults(targetDirectory, "../kibana/.github/CODEOWNERS", "color_palette_by_code_owner/results.csv", patterns)
}