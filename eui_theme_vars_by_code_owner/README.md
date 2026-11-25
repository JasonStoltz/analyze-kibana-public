# eui_theme_vars_by_code_owner

This script prints out the number of specific **EUI theme variable usages** in all files per directory under the Kibana repository.

It also prints the overall code-owner for that directory. Some directories do not have a code-owner, and it doesn't account for code-owners in nested files and directories, so this information can be of limited usefulness.

The script uses goroutines and a worker pool to process files concurrently. This is important for parsing big projects like Kibana.

## Patterns

By default, it searches for the following pattern:

- `euiThemeVars`

This can be changed by modifying the script.

## Usage

### Prerequisites

1. Clone Kibana into `../` folder.
2. Checkout the `main` branch.
3. `git pull` / `git pull upstream` (if using a forked repo).

### Steps

1. Go to the root folder: `cd analyze-kibana`.
2. Run `go run eui_theme_vars_by_code_owner/main.go` (pass the target directory as the first argument; `../kibana` by default).

You will see `results.csv` in the script's folder and the TSV will be copied to clipboard (to paste into a spreadsheet).
