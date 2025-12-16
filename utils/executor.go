package utils

import (
	"bytes"
	"os/exec"
)

type CommandResult struct {
	Command string `json:"command"`
	Output  string `json:"output"`
	Success bool   `json:"success"`
}

func Execute(cmd string) CommandResult {
	c := exec.Command("bash", "-c", cmd)

	var stdout, stderr bytes.Buffer
	c.Stdout = &stdout
	c.Stderr = &stderr

	err := c.Run()
	out := stdout.String() + stderr.String()

	return CommandResult{
		Command: cmd,
		Output:  out,
		Success: err == nil,
	}
}
