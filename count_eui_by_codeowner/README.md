This script prints out the number of EUI usages in all files in the `/plugins` directory in Kibana.

It counts usages in all files and groups the overall count per codeowner.

In addition to counts of EUI usage, it also counts the total number of files per codeowner and total number of .tsx files.

NOTE: This file was written in an incredibly ineffecient way. It walks the entire directory tree and greps each file one by one. A much more effective way would be to run the recursive grep once on the entire directory and then match only that subset of files that has matches to it's corresponding codeowner. As it is now, it takes around 5 minutes to run.

By default, it walks the `../kibana` directory looking for the text `<EUI`. To change that, pass a positional arguments with another directory and/or text to find:

```
go run count_eui_by_codeowner/main.go ../kibana/x-pack/platform/plugins/shared/stack_alerts "<button"
```
