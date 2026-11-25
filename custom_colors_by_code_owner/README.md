# custom_colors_by_code_owner

This script prints out the occurrences of custom color usages in all files per directory under the Kibana repository.

The script uses goroutines and a worker pool to process files concurrently. This is important for parsing big projects like Kibana.

## Patterns

By default, it searches for the following patterns:

- Hex colors (e.g., `#ff0000`, `#f00`)
- RGB colors (e.g., `rgb(255, 0, 0)`)
- RGBA colors (e.g., `rgba(255, 0, 0, 0.5)`)

This can be changed by modifying the script.

## Usage

### Prerequisites

1. Clone Kibana into `../` folder.
2. Checkout the `main` branch.
3. `git pull` / `git pull upstream` (if using a forked repo).

### Steps

1. Go to the root folder: `cd analyze-kibana`.
2. Run `go run custom_colors_by_code_owner/main.go` (pass the target directory as the first argument; `../kibana` by default).

You will see `results.csv` in the script's folder and the TSV will be copied to clipboard (to paste into a spreadsheet).
