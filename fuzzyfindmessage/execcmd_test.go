package fuzzyfindmessage

import (
	"fmt"
	"os/exec"
	"testing"
)

func Test__gitCommit(t *testing.T) {
	tests := []struct {
		name        string
		execCommand func(name string, arg ...string) *exec.Cmd
		commandRun  func(c *exec.Cmd) error
		fileName    string
		wantErr     bool
	}{
		{
			name: "Normal",
			execCommand: func(name string, arg ...string) *exec.Cmd {
				return &exec.Cmd{}
			},
			commandRun: func(c *exec.Cmd) error {
				return nil
			},
			fileName: "hoge",
			wantErr:  false,
		},
		{
			name: "ErrorBecauseCommandReturnError",
			execCommand: func(name string, arg ...string) *exec.Cmd {
				return &exec.Cmd{}
			},
			commandRun: func(c *exec.Cmd) error {
				return fmt.Errorf("error")
			},
			fileName: "hoge",
			wantErr:  true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execCommand = tt.execCommand
			commandRun = tt.commandRun
			if err := _gitCommit(tt.fileName); (err != nil) != tt.wantErr {
				t.Errorf("gitCommit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test__lastCommitMessage(t *testing.T) {
	tests := []struct {
		name          string
		execCommand   func(name string, arg ...string) *exec.Cmd
		commandOutput func(c *exec.Cmd) ([]byte, error)
		want          string
		wantErr       bool
	}{
		{
			name: "Normal",
			execCommand: func(name string, arg ...string) *exec.Cmd {
				return &exec.Cmd{}
			},
			commandOutput: func(c *exec.Cmd) ([]byte, error) {
				return []byte("hoge"), nil
			},
			want:    "hoge",
			wantErr: false,
		},
		{
			name: "ErrorBecauseCommandReturnError",
			execCommand: func(name string, arg ...string) *exec.Cmd {
				return &exec.Cmd{}
			},
			commandOutput: func(c *exec.Cmd) ([]byte, error) {
				return []byte(""), fmt.Errorf("error")
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execCommand = tt.execCommand
			commandOutput = tt.commandOutput
			got, err := _lastCommitMessage()
			if (err != nil) != tt.wantErr {
				t.Errorf("lastCommitMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("lastCommitMessage() got = %v, want %v", got, tt.want)
			}
		})
	}
}
