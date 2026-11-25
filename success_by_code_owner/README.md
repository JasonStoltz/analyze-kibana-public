# count_success_by_code_owner

This script prints out the number of specific **success color usages** in all files per directory under the Kibana repository.

It also prints the overall code-owner for that directory. Some directories do not have a code-owner, and it doesn't account for code-owners in nested files and directories, so this information can be of limited usefulness.

The script uses goroutines and a worker pool to process files concurrently. This is important for parsing big projects like Kibana.

## Patterns

By default, it searches for the following patterns:

- `color="success"`
- `$euiColorSuccess`
- `$euiColorSuccessText`
- `euiTheme.colors.success`
- `euiTheme.colors.successText`
- `theme.eui.euiColorSuccess`
- `theme.eui.euiColorSuccessText`

This can be changed by modifying the script.

## Usage

To run the script, use the following command:
`go run success_by_code_owner/main.go [target_directory]`

If no target directory is specified, it defaults to `../kibana`.

The output will display the code-owner along with the count of each pattern found in their files.
