package fuzzyfindmessage

import (
	"os"
	"os/exec"
)

var (
	execCommand = exec.Command
	commandRun  = func(c *exec.Cmd) error {
		return c.Run()
	}
	commandOutput = func(c *exec.Cmd) ([]byte, error) {
		return c.Output()
	}
)

func _gitCommit(fileName string) error {
	c := execCommand("git", "commit", "-F", fileName, "-e")
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return commandRun(c)
}

func _lastCommitMessage() (string, error) {
	c := execCommand("git", "log", "-1", "--pretty='%B'")
	out, err := commandOutput(c)
	if err != nil {
		return "", err
	}
	return string(out), nil
}
