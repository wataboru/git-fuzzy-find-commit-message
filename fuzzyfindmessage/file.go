// Package fuzzyfindmessage provides a Git commit message template fuzzy search feature.
// Optionally, it provides a fuzzy search for the history of your commit messages.
package fuzzyfindmessage

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/user"
	"sort"
	"strings"
	"time"

	"github.com/ktr0731/go-fuzzyfinder"
)

const (
	exampleFile = ".fcm"
	historyFile = ".fcm_history"
)

var (
	osOpen               func(name string) (*os.File, error)
	osOpenFile           func(name string, flag int, perm os.FileMode) (*os.File, error)
	osCreate             func(name string) (*os.File, error)
	osStat               func(name string) (os.FileInfo, error)
	osRemove             func(name string) error
	bufioNewScanner      func(r io.Reader) *bufio.Scanner
	fmtFprintf           func(w io.Writer, format string, a ...interface{}) (n int, err error)
	fmtFprintln          func(w io.Writer, a ...interface{}) (n int, err error)
	ioutilTempFile       func(dir, pattern string) (f *os.File, err error)
	fileWrite            func(file *os.File, b []byte) (n int, err error)
	fileClose            func(file *os.File) error
	scannerScan          func(scanner *bufio.Scanner) bool
	scannerText          func(scanner *bufio.Scanner) string
	newCurrentUser       func() *user.User
	home                 string
	userCurrent          func() (*user.User, error)
	exampleFilePath      string
	historyFilePath      string
	samples              func() ([]string, error)
	saveHistory          func() (err error)
	createTemplate       func(message string) (f *os.File, err error)
	createDefaultFile    func(filePath string) error
	removeDuplicate      func(slice []string) []string
	exists               func(filename string) bool
	createEmptyHistory   func() (err error)
	createDefaultExample func() (err error)
	lastCommitMessage    func() (string, error)
	gitCommit            func(fileName string) error
	fuzzyfinderFind      func(slice interface{}, itemFunc func(i int) string, opts ...fuzzyfinder.Option) (int, error)
	tmpFileName          func(f *os.File) string
)

func init() {
	osOpen = os.Open
	osOpenFile = os.OpenFile
	osCreate = os.Create
	osStat = os.Stat
	osRemove = os.Remove
	bufioNewScanner = bufio.NewScanner
	fmtFprintf = fmt.Fprintf
	fmtFprintln = fmt.Fprintln
	ioutilTempFile = ioutil.TempFile
	fileWrite = func(file *os.File, b []byte) (n int, err error) {
		return file.Write(b)
	}
	fileClose = func(file *os.File) error {
		return file.Close()
	}
	scannerScan = func(scanner *bufio.Scanner) bool {
		return scanner.Scan()
	}
	scannerText = func(scanner *bufio.Scanner) string {
		return scanner.Text()
	}
	userCurrent = user.Current
	newCurrentUser = _newCurrentUser
	home = newCurrentUser().HomeDir
	exampleFilePath = home + "/" + exampleFile
	historyFilePath = home + "/" + historyFile
	samples = _samples
	saveHistory = _saveHistory
	createTemplate = _createTemplate
	createDefaultFile = _createDefaultFile
	removeDuplicate = _removeDuplicate
	exists = _exists
	createEmptyHistory = _createEmptyHistory
	createDefaultExample = _createDefaultExample
	lastCommitMessage = _lastCommitMessage
	gitCommit = _gitCommit
	fuzzyfinderFind = fuzzyfinder.Find
	tmpFileName = func(f *os.File) string {
		return f.Name()
	}
}

// Commit wraps Git Commit.
// You can perform a fuzzy search from a message template and commit the result.
func Commit() (err error) {
	samples, err := samples()
	if err != nil {
		return err
	}

	id, err := fuzzyfinderFind(
		samples,
		func(i int) string {
			return samples[i]
		},
		fuzzyfinder.WithPreviewWindow(func(i, w, h int) string {
			if i == -1 {
				return ""
			}
			return fmt.Sprintln(samples[i])
		}))
	if err != nil {
		return err
	}

	f, err := createTemplate(samples[id])
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = osRemove(tmpFileName(f))
		}
		osRemove(tmpFileName(f))
	}()

	if err := gitCommit(tmpFileName(f)); err != nil {
		return err
	}

	if err := saveHistory(); err != nil {
		return err
	}

	return nil
}

func _createTemplate(message string) (f *os.File, err error) {
	f, err = ioutilTempFile("", "template")
	if err != nil {
		return nil, err
	}
	defer func() {
		if err == nil {
			err = fileClose(f)
		}
		fileClose(f)
	}()

	message = strings.ReplaceAll(message, "\\n", "\n")

	if _, err := fileWrite(f, []byte(message)); err != nil {
		return nil, err
	}

	return f, nil
}

func _samples() ([]string, error) {
	if err := createDefaultFile(exampleFilePath); err != nil {
		return nil, err
	}

	if err := createDefaultFile(historyFilePath); err != nil {
		return nil, err
	}

	var samples []string
	samplesFile, err := osOpen(exampleFilePath)
	if err != nil {
		return nil, err
	}
	defer fileClose(samplesFile)
	historyFile, err := osOpen(historyFilePath)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err == nil {
			err = fileClose(historyFile)
			return
		}
		fileClose(historyFile)
	}()

	scanner := bufioNewScanner(io.MultiReader(samplesFile, historyFile))
	for scannerScan(scanner) {
		s := scannerText(scanner)
		if len(s) == 0 {
			continue
		}

		if s[0:1] == "#" {
			continue
		}
		samples = append(samples, s)
	}

	samples = removeDuplicate(samples)

	sort.Slice(samples, func(i, j int) bool {
		return samples[i] > samples[j]
	})

	return samples, nil
}

func _saveHistory() (err error) {
	history, err := lastCommitMessage()
	if err != nil {
		return err
	}

	file, err := osOpenFile(historyFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = fileClose(file)
			return
		}
		fileClose(file)
	}()

	history = strings.Replace(history, "\n", "\\n", -1)

	if _, err := fmtFprintf(file, "# %s\n%s\n", time.Now().Format("2006/01/02 15:04:05"), history); err != nil {
		return err
	}

	return nil
}

func _removeDuplicate(slice []string) []string {
	results := make([]string, 0, len(slice))
	m := map[string]bool{}
	for i := range slice {
		if !m[slice[i]] {
			m[slice[i]] = true
			results = append(results, slice[i])
		}
	}
	return results
}

func _createDefaultFile(filePath string) error {
	if exists(filePath) {
		return nil
	}

	switch filePath {
	case historyFilePath:
		return createEmptyHistory()
	case exampleFilePath:
		return createDefaultExample()
	}

	return nil
}

func _createDefaultExample() (err error) {
	file, err := osCreate(exampleFilePath)
	if err != nil {
		return err
	}
	defer func() {
		if err == nil {
			err = fileClose(file)
			return
		}
		fileClose(file)
	}()

	for _, s := range defaultExamples {
		if _, err := fmtFprintln(file, s); err != nil {
			return err
		}
	}

	return nil
}

func _createEmptyHistory() (err error) {
	file, err := osCreate(historyFilePath)
	if err != nil {
		return err
	}
	defer func() {
		err = fileClose(file)
	}()

	return nil
}

func _exists(filename string) bool {
	_, err := osStat(filename)
	return err == nil
}

func _newCurrentUser() *user.User {
	usr, err := userCurrent()
	if err != nil {
		panic("error")
	}
	return usr
}
