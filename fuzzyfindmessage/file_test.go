package fuzzyfindmessage

import (
	"bufio"
	"fmt"
	"github.com/ktr0731/go-fuzzyfinder"
	"io"
	"os"
	"os/user"
	"reflect"
	"testing"
)

func Test_samples(t *testing.T) {
	count := 0
	tests := []struct {
		name              string
		createDefaultFile func(filePath string) error
		osOpen            func(name string) (*os.File, error)
		bufioNewScanner   func(r io.Reader) *bufio.Scanner
		scannerScan       func(scanner *bufio.Scanner) bool
		scannerText       func(scanner *bufio.Scanner) string
		removeDuplicate   func(slice []string) []string
		want              []string
		wantErr           bool
	}{
		{
			name: "Normal",
			createDefaultFile: func(filePath string) error {
				return nil
			},
			osOpen: func(name string) (*os.File, error) {
				return nil, nil
			},
			bufioNewScanner: func(r io.Reader) *bufio.Scanner {
				return &bufio.Scanner{}
			},
			scannerScan: func(scanner *bufio.Scanner) bool {
				count++
				return count <= 2
			},
			scannerText: func(scanner *bufio.Scanner) string {
				return "hoge"
			},
			removeDuplicate: func(slice []string) []string {
				return []string{"fuga", "hoge"}
			},
			want:    []string{"hoge", "fuga"},
			wantErr: false,
		},
		{
			name: "NormalNoneRecordInFile",
			createDefaultFile: func(filePath string) error {
				return nil
			},
			osOpen: func(name string) (*os.File, error) {
				return nil, nil
			},
			bufioNewScanner: func(r io.Reader) *bufio.Scanner {
				return &bufio.Scanner{}
			},
			scannerScan: func(scanner *bufio.Scanner) bool {
				return false
			},
			scannerText: nil,
			removeDuplicate: func(slice []string) []string {
				return nil
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "NormalBlankLine",
			createDefaultFile: func(filePath string) error {
				return nil
			},
			osOpen: func(name string) (*os.File, error) {
				return nil, nil
			},
			bufioNewScanner: func(r io.Reader) *bufio.Scanner {
				return &bufio.Scanner{}
			},
			scannerScan: func(scanner *bufio.Scanner) bool {
				count++
				return count <= 1

			},
			scannerText: func(scanner *bufio.Scanner) string {
				return ""
			},
			removeDuplicate: func(slice []string) []string {
				return nil
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "NormalBlankLine",
			createDefaultFile: func(filePath string) error {
				return nil
			},
			osOpen: func(name string) (*os.File, error) {
				return nil, nil
			},
			bufioNewScanner: func(r io.Reader) *bufio.Scanner {
				return &bufio.Scanner{}
			},
			scannerScan: func(scanner *bufio.Scanner) bool {
				count++
				return count <= 1
			},
			scannerText: func(scanner *bufio.Scanner) string {
				return "# hoge"
			},
			removeDuplicate: func(slice []string) []string {
				return nil
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "ErrorBecauseCreateExamplesDefaultFileError",
			createDefaultFile: func(filePath string) error {
				return fmt.Errorf("error")
			},
			osOpen: nil,
			bufioNewScanner: func(r io.Reader) *bufio.Scanner {
				return &bufio.Scanner{}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ErrorBecauseCreateHistoryDefaultFileError",
			createDefaultFile: func(filePath string) error {
				count++
				if count <= 1 {
					return nil
				}
				return fmt.Errorf("error")
			},
			osOpen: nil,
			bufioNewScanner: func(r io.Reader) *bufio.Scanner {
				return &bufio.Scanner{}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ErrorBecauseNotOpenExamplesFile",
			createDefaultFile: func(filePath string) error {
				return nil
			},
			osOpen: func(name string) (*os.File, error) {
				return nil, fmt.Errorf("error")
			},
			bufioNewScanner: func(r io.Reader) *bufio.Scanner {
				return &bufio.Scanner{}
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "ErrorBecauseNotOpenHistoryFile",
			createDefaultFile: func(filePath string) error {
				return nil
			},
			osOpen: func(name string) (*os.File, error) {
				count++
				if count <= 1 {
					return nil, nil
				}
				return nil, fmt.Errorf("error")
			},
			bufioNewScanner: func(r io.Reader) *bufio.Scanner {
				return &bufio.Scanner{}
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			count = 0
			createDefaultFile = tt.createDefaultFile
			osOpen = tt.osOpen
			bufioNewScanner = tt.bufioNewScanner
			scannerScan = tt.scannerScan
			scannerText = tt.scannerText
			removeDuplicate = tt.removeDuplicate
			got, err := _samples()
			if (err != nil) != tt.wantErr {
				t.Errorf("samples() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("samples() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test__saveHistory(t *testing.T) {
	tests := []struct {
		name              string
		lastCommitMessage func() (string, error)
		osOpenFile        func(name string, flag int, perm os.FileMode) (*os.File, error)
		fileClose         func(file *os.File) error
		fmtFprintf        func(w io.Writer, format string, a ...interface{}) (n int, err error)
		wantErr           bool
	}{
		{
			name: "Normal",
			lastCommitMessage: func() (string, error) {
				return "", nil
			},
			osOpenFile: func(name string, flag int, perm os.FileMode) (*os.File, error) {
				return nil, nil
			},
			fileClose: func(file *os.File) error {
				return nil
			},
			fmtFprintf: func(w io.Writer, format string, a ...interface{}) (n int, err error) {
				return n, nil
			},
			wantErr: false,
		},
		{
			name: "ErrorLastCommitMessageReturnError",
			lastCommitMessage: func() (string, error) {
				return "", fmt.Errorf("error")
			},
			osOpenFile: func(name string, flag int, perm os.FileMode) (*os.File, error) {
				return nil, fmt.Errorf("error")
			},
			fileClose: func(file *os.File) error {
				return nil
			},
			fmtFprintf: func(w io.Writer, format string, a ...interface{}) (n int, err error) {
				return n, nil
			},
			wantErr: true,
		},
		{
			name: "ErrorBecauseNotOpenFile",
			lastCommitMessage: func() (string, error) {
				return "", nil
			},
			osOpenFile: func(name string, flag int, perm os.FileMode) (*os.File, error) {
				return nil, fmt.Errorf("error")
			},
			fileClose: func(file *os.File) error {
				return nil
			},
			fmtFprintf: func(w io.Writer, format string, a ...interface{}) (n int, err error) {
				return n, nil
			},
			wantErr: true,
		},
		{
			name: "ErrorBecauseFprintReturnError",
			lastCommitMessage: func() (string, error) {
				return "", nil
			},
			osOpenFile: func(name string, flag int, perm os.FileMode) (*os.File, error) {
				return nil, nil
			},
			fileClose: func(file *os.File) error {
				return nil
			},
			fmtFprintf: func(w io.Writer, format string, a ...interface{}) (n int, err error) {
				return n, fmt.Errorf("error")
			},
			wantErr: true,
		},
		{
			name: "ErrorBecauseFileCloseReturnError",
			lastCommitMessage: func() (string, error) {
				return "", nil
			},
			osOpenFile: func(name string, flag int, perm os.FileMode) (*os.File, error) {
				return nil, nil
			},
			fileClose: func(file *os.File) error {
				return fmt.Errorf("error")
			},
			fmtFprintf: func(w io.Writer, format string, a ...interface{}) (n int, err error) {
				return
			},
			wantErr: true,
		},
		{
			name: "ErrorBecauseFileCloseReturnErrorAndFileCloseReturnError",
			lastCommitMessage: func() (string, error) {
				return "", nil
			},
			osOpenFile: func(name string, flag int, perm os.FileMode) (*os.File, error) {
				return nil, nil
			},
			fileClose: func(file *os.File) error {
				return fmt.Errorf("error")
			},
			fmtFprintf: func(w io.Writer, format string, a ...interface{}) (n int, err error) {
				return n, fmt.Errorf("error")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lastCommitMessage = tt.lastCommitMessage
			osOpenFile = tt.osOpenFile
			fileClose = tt.fileClose
			fmtFprintf = tt.fmtFprintf
			if err := _saveHistory(); (err != nil) != tt.wantErr {
				t.Errorf("_saveHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test__newCurrentUser(t *testing.T) {
	tests := []struct {
		name        string
		userCurrent func() (*user.User, error)
		want        *user.User
		wantErr     bool
	}{
		{
			name: "Normal",
			userCurrent: func() (*user.User, error) {
				return nil, nil
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "PanicBecauseReturnUserCurrentReturnError",
			userCurrent: func() (*user.User, error) {
				return nil, fmt.Errorf("error")
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userCurrent = tt.userCurrent
			if tt.wantErr {
				defer func() {
					err := recover()
					if err == nil {
						t.Errorf("newCurrentUser() not return panic")
					}
				}()
			}
			if got := _newCurrentUser(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("newCurrentUser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test__exists(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		osStat   func(name string) (os.FileInfo, error)
		want     bool
	}{
		{
			name:     "True",
			filename: "hoge",
			osStat: func(name string) (os.FileInfo, error) {
				return nil, nil
			},
			want: true,
		},
		{
			name:     "False",
			filename: "hoge",
			osStat: func(name string) (os.FileInfo, error) {
				return nil, fmt.Errorf("error")
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			osStat = tt.osStat
			if got := _exists(tt.filename); got != tt.want {
				t.Errorf("exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test__createEmptyHistory(t *testing.T) {
	tests := []struct {
		name      string
		osCreate  func(name string) (*os.File, error)
		fileClose func(file *os.File) error
		wantErr   bool
	}{
		{
			name: "Normal",
			osCreate: func(name string) (*os.File, error) {
				return nil, nil
			},
			fileClose: func(file *os.File) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "ErrorBecauseOsCreateReturnError",
			osCreate: func(name string) (*os.File, error) {
				return nil, fmt.Errorf("error")
			},
			fileClose: func(file *os.File) error {
				return nil
			},
			wantErr: true,
		},
		{
			name: "ErrorBecauseFileCloseReturnError",
			osCreate: func(name string) (*os.File, error) {
				return nil, nil
			},
			fileClose: func(file *os.File) error {
				return fmt.Errorf("error")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			osCreate = tt.osCreate
			fileClose = tt.fileClose
			if err := _createEmptyHistory(); (err != nil) != tt.wantErr {
				t.Errorf("createEmptyHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test__createDefaultExample(t *testing.T) {
	tests := []struct {
		name        string
		osCreate    func(name string) (*os.File, error)
		fileClose   func(file *os.File) error
		fmtFprintln func(w io.Writer, a ...interface{}) (n int, err error)
		wantErr     bool
	}{
		{
			name: "Normal",
			osCreate: func(name string) (*os.File, error) {
				return nil, nil
			},
			fileClose: func(file *os.File) error {
				return nil
			},
			fmtFprintln: func(w io.Writer, a ...interface{}) (n int, err error) {
				return
			},
			wantErr: false,
		},
		{
			name: "ErrorBecauseOsCreateReturnError",
			osCreate: func(name string) (*os.File, error) {
				return nil, fmt.Errorf("error")
			},
			fileClose: func(file *os.File) error {
				return nil
			},
			fmtFprintln: func(w io.Writer, a ...interface{}) (n int, err error) {
				return
			},
			wantErr: true,
		},
		{
			name: "ErrorBecauseFmtFprintlnReturnError",
			osCreate: func(name string) (*os.File, error) {
				return nil, nil
			},
			fileClose: func(file *os.File) error {
				return nil
			},
			fmtFprintln: func(w io.Writer, a ...interface{}) (n int, err error) {
				return n, fmt.Errorf("error")
			},
			wantErr: true,
		},
		{
			name: "ErrorBecauseFileCloseReturnError",
			osCreate: func(name string) (*os.File, error) {
				return nil, nil
			},
			fileClose: func(file *os.File) error {
				return fmt.Errorf("error")
			},
			fmtFprintln: func(w io.Writer, a ...interface{}) (n int, err error) {
				return
			},
			wantErr: true,
		},
		{
			name: "ErrorBecauseFileCloseReturnError",
			osCreate: func(name string) (*os.File, error) {
				return nil, nil
			},
			fileClose: func(file *os.File) error {
				return fmt.Errorf("error")
			},
			fmtFprintln: func(w io.Writer, a ...interface{}) (n int, err error) {
				return n, fmt.Errorf("error")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			osCreate = tt.osCreate
			fileClose = tt.fileClose
			fmtFprintln = tt.fmtFprintln
			if err := _createDefaultExample(); (err != nil) != tt.wantErr {
				t.Errorf("createDefaultExample() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test__createTemplate(t *testing.T) {
	tests := []struct {
		name           string
		message        string
		ioutilTempFile func(dir, pattern string) (f *os.File, err error)
		fileWrite      func(file *os.File, b []byte) (n int, err error)
		fileClose      func(file *os.File) error
		want           *os.File
		wantErr        bool
	}{
		{
			name:    "Normal",
			message: "hoge",
			ioutilTempFile: func(dir, pattern string) (f *os.File, err error) {
				return
			},
			fileWrite: func(file *os.File, b []byte) (n int, err error) {
				return

			},
			fileClose: func(file *os.File) error {
				return nil
			},
			want:    nil,
			wantErr: false,
		},
		{
			name:    "ErrorBecauseIoutilTempFileReturnError",
			message: "hoge",
			ioutilTempFile: func(dir, pattern string) (f *os.File, err error) {
				return nil, fmt.Errorf("error")
			},
			fileWrite: func(file *os.File, b []byte) (n int, err error) {
				return

			},
			fileClose: func(file *os.File) error {
				return nil
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ErrorBecauseFileWriteReturnError",
			message: "hoge",
			ioutilTempFile: func(dir, pattern string) (f *os.File, err error) {
				return
			},
			fileWrite: func(file *os.File, b []byte) (n int, err error) {
				return n, fmt.Errorf("error")

			},
			fileClose: func(file *os.File) error {
				return nil
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ErrorBecauseFileCloseReturnError",
			message: "hoge",
			ioutilTempFile: func(dir, pattern string) (f *os.File, err error) {
				return
			},
			fileWrite: func(file *os.File, b []byte) (n int, err error) {
				return

			},
			fileClose: func(file *os.File) error {
				return fmt.Errorf("error")
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "ErrorBecauseFileWriteReturnErrorAndFileCloseReturnError",
			message: "hoge",
			ioutilTempFile: func(dir, pattern string) (f *os.File, err error) {
				return
			},
			fileWrite: func(file *os.File, b []byte) (n int, err error) {
				return n, fmt.Errorf("error")

			},
			fileClose: func(file *os.File) error {
				return fmt.Errorf("error")
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ioutilTempFile = tt.ioutilTempFile
			fileWrite = tt.fileWrite
			fileClose = tt.fileClose
			got, err := _createTemplate(tt.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("createTemplate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("createTemplate() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test__createDefaultFile(t *testing.T) {
	historyFilePathMock := "hoge/.fcm_history"
	exampleFilePathMock := "hoge/.fcm"
	tests := []struct {
		name                 string
		filePath             string
		exists               func(filename string) bool
		createEmptyHistory   func() (err error)
		createDefaultExample func() (err error)
		wantErr              bool
	}{
		{
			name:     "FileExist",
			filePath: "hoge",
			exists: func(filename string) bool {
				return true
			},
			createEmptyHistory: func() (err error) {
				return nil
			},
			createDefaultExample: func() (err error) {
				return nil
			},
			wantErr: false,
		},
		{
			name:     "createEmptyHistoryNotReturnError",
			filePath: "hoge/.fcm_history",
			exists: func(filename string) bool {
				return false
			},
			createEmptyHistory: func() (err error) {
				return nil
			},
			createDefaultExample: func() (err error) {
				return nil
			},
			wantErr: false,
		},
		{
			name:     "createDefaultExampleNotReturnError",
			filePath: "hoge/.fcm",
			exists: func(filename string) bool {
				return false
			},
			createEmptyHistory: func() (err error) {
				return nil
			},
			createDefaultExample: func() (err error) {
				return nil
			},
			wantErr: false,
		},
		{
			name:     "nonTargetFilepath",
			filePath: "hoge/hoge.txt",
			exists: func(filename string) bool {
				return false
			},
			createEmptyHistory:   nil,
			createDefaultExample: nil,
			wantErr:              false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			exists = tt.exists
			createEmptyHistory = tt.createEmptyHistory
			createDefaultExample = tt.createDefaultExample
			exampleFilePath = exampleFilePathMock
			historyFilePath = historyFilePathMock
			if err := _createDefaultFile(tt.filePath); (err != nil) != tt.wantErr {
				t.Errorf("_createDefaultFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test__removeDuplicate(t *testing.T) {
	tests := []struct {
		name  string
		slice []string
		want  []string
	}{
		{
			name:  "Normal",
			slice: []string{"hoge", "hoge"},
			want:  []string{"hoge"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := _removeDuplicate(tt.slice); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("_removeDuplicate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommit(t *testing.T) {
	tests := []struct {
		name            string
		samples         func() ([]string, error)
		fuzzyfinderFind func(slice interface{}, itemFunc func(i int) string, opts ...fuzzyfinder.Option) (int, error)
		createTemplate  func(message string) (f *os.File, err error)
		gitCommit       func(fileName string) error
		saveHistory     func() (err error)
		tmpFileName     func(f *os.File) string
		osRemove        func(name string) error
		wantErr         bool
	}{
		{
			name: "Normal",
			samples: func() ([]string, error) {
				return []string{"hoge"}, nil
			},
			fuzzyfinderFind: func(slice interface{}, itemFunc func(i int) string, opts ...fuzzyfinder.Option) (int, error) {
				return 0, nil
			},
			createTemplate: func(message string) (f *os.File, err error) {
				return nil, nil
			},
			gitCommit: func(fileName string) error {
				return nil
			},
			saveHistory: func() (err error) {
				return nil
			},
			tmpFileName: func(f *os.File) string {
				return "hoge"
			},
			osRemove: func(name string) error {
				return nil
			},
			wantErr: false,
		},
		{
			name: "ErrorBecauseSamplesReturnError",
			samples: func() ([]string, error) {
				return nil, fmt.Errorf("error")
			},
			fuzzyfinderFind: func(slice interface{}, itemFunc func(i int) string, opts ...fuzzyfinder.Option) (int, error) {
				return 0, nil
			},
			createTemplate: func(message string) (f *os.File, err error) {
				return nil, nil
			},
			gitCommit: func(fileName string) error {
				return nil
			},
			saveHistory: func() (err error) {
				return nil
			},
			tmpFileName: func(f *os.File) string {
				return "hoge"
			},
			osRemove: func(name string) error {
				return nil
			},
			wantErr: true,
		},
		{
			name: "ErrorBecauseFuzzyFinderFindReturnError",
			samples: func() ([]string, error) {
				return []string{"hoge"}, nil
			},
			fuzzyfinderFind: func(slice interface{}, itemFunc func(i int) string, opts ...fuzzyfinder.Option) (int, error) {
				return 0, fmt.Errorf("error")
			},
			createTemplate: func(message string) (f *os.File, err error) {
				return nil, nil
			},
			gitCommit: func(fileName string) error {
				return nil
			},
			saveHistory: func() (err error) {
				return nil
			},
			tmpFileName: func(f *os.File) string {
				return "hoge"
			},
			osRemove: func(name string) error {
				return nil
			},
			wantErr: true,
		},
		{
			name: "ErrorBecauseCreateTemplateReturnError",
			samples: func() ([]string, error) {
				return []string{"hoge"}, nil
			},
			fuzzyfinderFind: func(slice interface{}, itemFunc func(i int) string, opts ...fuzzyfinder.Option) (int, error) {
				return 0, nil
			},
			createTemplate: func(message string) (f *os.File, err error) {
				return nil, fmt.Errorf("error")
			},
			gitCommit: func(fileName string) error {
				return nil
			},
			saveHistory: func() (err error) {
				return nil
			},
			tmpFileName: func(f *os.File) string {
				return "hoge"
			},
			osRemove: func(name string) error {
				return nil
			},
			wantErr: true,
		},
		{
			name: "ErrorBecauseGitCommitReturnError",
			samples: func() ([]string, error) {
				return []string{"hoge"}, nil
			},
			fuzzyfinderFind: func(slice interface{}, itemFunc func(i int) string, opts ...fuzzyfinder.Option) (int, error) {
				return 0, nil
			},
			createTemplate: func(message string) (f *os.File, err error) {
				return nil, nil
			},
			gitCommit: func(fileName string) error {
				return fmt.Errorf("error")
			},
			saveHistory: func() (err error) {
				return nil
			},
			tmpFileName: func(f *os.File) string {
				return "hoge"
			},
			osRemove: func(name string) error {
				return nil
			},
			wantErr: true,
		},
		{
			name: "ErrorBecauseSaveHistoryReturnError",
			samples: func() ([]string, error) {
				return []string{"hoge"}, nil
			},
			fuzzyfinderFind: func(slice interface{}, itemFunc func(i int) string, opts ...fuzzyfinder.Option) (int, error) {
				return 0, nil
			},
			createTemplate: func(message string) (f *os.File, err error) {
				return nil, nil
			},
			gitCommit: func(fileName string) error {
				return nil
			},
			saveHistory: func() (err error) {
				return fmt.Errorf("error")
			},
			tmpFileName: func(f *os.File) string {
				return "hoge"
			},
			osRemove: func(name string) error {
				return nil
			},
			wantErr: true,
		},
		{
			name: "ErrorBecauseOsRemoveReturnError",
			samples: func() ([]string, error) {
				return []string{"hoge"}, nil
			},
			fuzzyfinderFind: func(slice interface{}, itemFunc func(i int) string, opts ...fuzzyfinder.Option) (int, error) {
				return 0, nil
			},
			createTemplate: func(message string) (f *os.File, err error) {
				return nil, nil
			},
			gitCommit: func(fileName string) error {
				return nil
			},
			saveHistory: func() (err error) {
				return nil
			},
			tmpFileName: func(f *os.File) string {
				return "hoge"
			},
			osRemove: func(name string) error {
				return fmt.Errorf("error")
			},
			wantErr: true,
		},
		{
			name: "ErrorBecauseOsRemoveReturnErrorAndSomeError",
			samples: func() ([]string, error) {
				return []string{"hoge"}, nil
			},
			fuzzyfinderFind: func(slice interface{}, itemFunc func(i int) string, opts ...fuzzyfinder.Option) (int, error) {
				return 0, nil
			},
			createTemplate: func(message string) (f *os.File, err error) {
				return nil, nil
			},
			gitCommit: func(fileName string) error {
				return fmt.Errorf("error")
			},
			saveHistory: func() (err error) {
				return nil
			},
			tmpFileName: func(f *os.File) string {
				return "hoge"
			},
			osRemove: func(name string) error {
				return fmt.Errorf("error")
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			samples = tt.samples
			fuzzyfinderFind = tt.fuzzyfinderFind
			createTemplate = tt.createTemplate
			gitCommit = tt.gitCommit
			saveHistory = tt.saveHistory
			tmpFileName = tt.tmpFileName
			osRemove = tt.osRemove
			if err := Commit(); (err != nil) != tt.wantErr {
				t.Errorf("Commit() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
