This code will read a list of Kibana filenames and output the corresponding codeowners.

They will print out in the same orders as the input, and include "null" as the value if a codeowner cannot be found.

Create `text-of-filenames-to-codeowner/tmp.txt` with one Kibana-relative path per line, then run the following command from the repository root.

```
go run text-of-filenames-to-codeowner/main.go
```
