# Analyze Kibana

These Go commands report how Kibana uses Elastic UI (EUI) and group findings by plugin or code owner.

## Run a report

After completing the fresh checkout steps below, run commands from the repository root. To list code owners, run

```sh
go run ./all_kibana_codeowners
```

To count matches for `<EuiButton` in each X-Pack plugin and save the tab-separated output:

```sh
go run ./count_eui_in_plugins '<EuiButton' > plugins.tsv
```

To analyze ownership and EUI usage in a single directory:

```sh
go run ./count_eui_by_codeowner ../kibana/x-pack/solutions/search/plugins/enterprise_search
```

The report directories contain additional usage examples. Color and theme reports write `results.csv` inside their own directories.

## Fresh checkout

1. Install [Go](https://go.dev/doc/install/) 1.22.2 or later and GNU grep with PCRE support under the name `ggrep`. On macOS with Homebrew, run `brew install go grep`.
2. From the root of this repository, run `git clone --depth 1 https://github.com/elastic/kibana.git ../kibana` if the sibling checkout is missing.
3. Run `go mod download` and `go build ./...`. With the tools on your shell's PATH, run the report commands above directly.

The optional `index_in_elasticsearch` command uses `ELASTIC_CLOUD_ID` and `ELASTIC_API_KEY` from the environment and writes to the `eui` index.
