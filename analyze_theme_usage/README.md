# Analyze theme usage

Count EUI component, theme variable, and custom color usages by code owner. Run the command from the repository root. An optional path limits the scan to part of the Kibana checkout.

```sh
go run ./analyze_theme_usage ../kibana/src
```

The command writes tab-separated results to standard output.
