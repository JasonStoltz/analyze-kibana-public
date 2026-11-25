package utils

import (
    "bufio"
    "encoding/csv"
    "fmt"
    "io"
    "io/fs"
    "log"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "sync"

    "github.com/atotto/clipboard"
    "github.com/gobwas/glob"
    "github.com/hmarr/codeowners"
)

var IncludeList = []glob.Glob{
    glob.MustCompile("*.ts"),
    glob.MustCompile("*.js"),
    glob.MustCompile("*.tsx"),
    glob.MustCompile("*.scss"),
    glob.MustCompile("*.yaml"),
    glob.MustCompile("*.html"),
}

var ExcludeList = []glob.Glob{
    glob.MustCompile("*examples/*"),
    glob.MustCompile("*.test*"),
    glob.MustCompile("*stories*"),
    glob.MustCompile("*.asciidoc"),
    glob.MustCompile("*.ndjson"),
    glob.MustCompile("*.mdx"),
    glob.MustCompile("*test*"),
    glob.MustCompile("*.md"),
    glob.MustCompile("*.svg"),
    glob.MustCompile("*.lock"),
}

func ProcessAndWriteResults(targetDirectory, codeownersPath, outputFile string, patterns map[string]*regexp.Regexp) {
    ruleset := loadCodeowners(codeownersPath)

    results := processFiles(targetDirectory, ruleset, patterns)

    header := []string{"Code-owner"}
    for pattern := range patterns {
        header = append(header, pattern)
    }

    writeResultsToCSV(results, outputFile, header)
}

/* Helpers */

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func grepCount(filePath string, pattern *regexp.Regexp) int {
    file, err := os.Open(filePath)
    check(err)
    defer file.Close()

    reader := bufio.NewReader(file)
    count := 0
    for {
        line, err := reader.ReadString('\n')
        if err != nil && err != io.EOF {
            check(err)
        }
        count += len(pattern.FindAllStringIndex(line, -1))
        if err == io.EOF {
            break
        }
    }

    return count
}

func worker(fileChan <-chan string, resultChan chan<- map[string]map[string]int, progressChan chan<- int, wg *sync.WaitGroup, ruleset codeowners.Ruleset, patterns map[string]*regexp.Regexp) {
    defer wg.Done()
    for file := range fileChan {
        rule, err := ruleset.Match(strings.Replace(file, "../kibana/", "", 1))
        check(err)

        if rule != nil {
            fileResults := make(map[string]map[string]int)
            for _, owner := range rule.Owners {
                ownerResults := make(map[string]int)
                for key, pattern := range patterns {
                    count := grepCount(file, pattern)
                    ownerResults[key] += count
                }
                fileResults[owner.Value] = ownerResults
            }
            resultChan <- fileResults
        }
        progressChan <- 1
    }
}

func printProgressBar(current, total int) {
    percent := float64(current) / float64(total) * 100
    barLength := 40
    filledLength := int(float64(barLength) * percent / 100)
    bar := strings.Repeat("=", filledLength) + strings.Repeat(" ", barLength-filledLength)
    fmt.Printf("\r[%s] %.2f%% (%d/%d)", bar, percent, current, total)

    if current >= total {
        fmt.Println()
    }
}

func matchAny(path string, patterns []glob.Glob) bool {
    for _, pattern := range patterns {
        if pattern.Match(path) {
            return true
        }
    }
    return false
}

func loadCodeowners(path string) codeowners.Ruleset {
    file, err := os.Open(path)
    check(err)
    defer file.Close()

    ruleset, err := codeowners.ParseFile(file)
    check(err)

    return ruleset
}

func processFiles(targetDirectory string, ruleset codeowners.Ruleset, patterns map[string]*regexp.Regexp) map[string]map[string]int {
    log.Println("Starting file processing...")
    results := make(map[string]map[string]int)
    var wg sync.WaitGroup
    resultChan := make(chan map[string]map[string]int, 100)
    progressChan := make(chan int, 100)

    numWorkers := 10
    fileChan := make(chan string, 100)

    totalFiles := 0
    filepath.WalkDir(targetDirectory, func(s string, d fs.DirEntry, err error) error {
        check(err)
        isExcluded := matchAny(s, ExcludeList)
        if !d.IsDir() && matchAny(s, IncludeList) && !isExcluded {
            totalFiles++
        }
        return nil
    })

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go worker(fileChan, resultChan, progressChan, &wg, ruleset, patterns)
    }

    go func() {
        filepath.WalkDir(targetDirectory, func(s string, d fs.DirEntry, err error) error {
            check(err)
            isExcluded := matchAny(s, ExcludeList)
            if !d.IsDir() && matchAny(s, IncludeList) && !isExcluded {
                fileChan <- s
            }
            return nil
        })
        close(fileChan)
    }()

    go func() {
        wg.Wait()
        close(resultChan)
        close(progressChan)
    }()

    processedFiles := 0

    go func() {
        for range progressChan {
            processedFiles++
            printProgressBar(processedFiles, totalFiles)
        }
    }()

    for fileResults := range resultChan {
        for owner, counts := range fileResults {
            if _, exists := results[owner]; !exists {
                results[owner] = make(map[string]int)
            }
            for key, count := range counts {
                results[owner][key] += count
            }
        }
    }

    log.Println("File processing completed.")
    return results
}

func writeResultsToCSV(results map[string]map[string]int, filename string, header []string) {
    file, err := os.Create(filename)
    check(err)
    defer file.Close()

    writer := csv.NewWriter(file)

    check(writer.Write(header))

    var clipboardContent strings.Builder
    clipboardContent.WriteString(strings.Join(header, "\t") + "\n")

    for owner, counts := range results {
        row := []string{owner}
        for _, key := range header[1:] {
            row = append(row, fmt.Sprint(counts[key]))
        }
        check(writer.Write(row))
        clipboardContent.WriteString(strings.Join(row, "\t") + "\n")
    }

    writer.Flush()
    check(writer.Error())
    log.Printf("Results written to %s.", filename)
    if err := clipboard.WriteAll(clipboardContent.String()); err != nil {
        log.Printf("Could not copy results to clipboard: %v", err)
    }
}
