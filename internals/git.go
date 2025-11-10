package internals

import (
	"bytes"
	"fmt"
	"os/exec"
)

func RunGit(args ...string) string {
	cmd := exec.Command("git", args...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()

	if err != nil {
		return fmt.Sprintf("Error: %v\n%s", err, out.String())
	}
	return out.String()
}
