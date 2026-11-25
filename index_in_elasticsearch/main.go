package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esutil"
	"github.com/hmarr/codeowners"
	"github.com/samber/lo"
)

//func UNUSED(x ...interface{}) {}

func main() {
	file, err := os.Open("../kibana/.github/CODEOWNERS")
	if err != nil {
		log.Fatal(err)
	}

	ruleset, err := codeowners.ParseFile(file)

	if err != nil {
		log.Fatal(err)
	}

	var numOfFiles int

	targetDirectory := "../kibana"
	if len(os.Args) > 1 {
		targetDirectory = os.Args[1]
	}

	cfg := elasticsearch.Config{
		CloudID: os.Getenv("ELASTIC_CLOUD_ID"),
		APIKey:  os.Getenv("ELASTIC_API_KEY"),
	}
	if cfg.CloudID == "" || cfg.APIKey == "" {
		log.Fatal("set ELASTIC_CLOUD_ID and ELASTIC_API_KEY to run the indexer")
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatal(err)
	}
	bi, err := esutil.NewBulkIndexer(esutil.BulkIndexerConfig{
		Index:         "eui",            // The default index name
		Client:        es,               // The Elasticsearch client
		NumWorkers:    runtime.NumCPU(), // The number of worker goroutines
		FlushBytes:    5e+6,             // The flush threshold in bytes
		FlushInterval: 30 * time.Second, // The periodic flush interval
	})
	if err != nil {
		log.Fatal(err)
	}

	timestamp := time.Now().UTC()

	err = filepath.WalkDir(targetDirectory, func(s string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Type().IsRegular() {
			numOfFiles++
			fmt.Printf("\rOn file %d", numOfFiles)

			rule, err := ruleset.Match(strings.Replace(s, "../kibana/", "", 1))
			if err != nil {
				println(err)
			}

			if rule != nil {
				ownersStringArr := lo.Map(rule.Owners, func(owner codeowners.Owner, index int) string {
					return owner.Value
				})

				out, err := exec.Command("grep", "-Ioh", "--", "<Eui[a-zA-Z0-9]*", s).Output()
				if err != nil {
					if exit, ok := err.(*exec.ExitError); !ok || exit.ExitCode() != 1 {
						return fmt.Errorf("grep %s: %w", s, err)
					}
				}
				for _, match := range strings.Fields(string(out)) {
					component := strings.TrimPrefix(match, "<")
					document := struct {
						FileName  string    `json:"filename"`
						Component string    `json:"component"`
						Timestamp time.Time `json:"@timestamp"`
						Codeowner []string  `json:"codeowners"`
						Filetype  string    `json:"filetype"`
					}{
						s,
						component,
						timestamp,
						ownersStringArr,
						filepath.Ext(s),
					}
					data, err := json.Marshal(document)
					if err != nil {
						log.Fatal(err)
					}
					err = bi.Add(
						context.Background(),
						esutil.BulkIndexerItem{
							Action: "index",
							Body:   bytes.NewReader(data),
						},
					)
					if err != nil {
						log.Fatalf("Unexpected error: %s", err)
					}
				}
			}
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	if err := bi.Close(context.Background()); err != nil {
		log.Fatalf("Unexpected error: %s", err)
	}
}
