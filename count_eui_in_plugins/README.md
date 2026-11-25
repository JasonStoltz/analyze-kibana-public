# Count EUI usage in X-Pack plugins

From the repository root, run:

```sh
go run ./count_eui_in_plugins '<EuiButton'
```

The optional argument is a GNU grep basic regular expression. It defaults to `Eui`. The tab-separated output contains each plugin name, its total number of matches, and the CODEOWNERS entry for its directory.

The command discovers plugins in the X-Pack platform and solution directories. It also supports older checkouts with an `x-pack/plugins` directory.
