package utils

import (
	"bytes"
	"fmt"
	"os/exec"
)

// CountGrepMatches counts the lines produced by GNU grep's only-matching mode.
func CountGrepMatches(options, pattern, path string) (int, error) {
	out, err := exec.Command("ggrep", options, "--", pattern, path).Output()
	if err != nil {
		if exit, ok := err.(*exec.ExitError); ok {
			if exit.ExitCode() == 1 {
				return 0, nil // grep found no matches
			}
			return 0, fmt.Errorf("ggrep %s: %w: %s", path, err, bytes.TrimSpace(exit.Stderr))
		}
		return 0, err
	}
	return bytes.Count(out, []byte("\n")), nil
}
