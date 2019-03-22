package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

const (
	cfgFile string = ".git-walk"
)

func moveTo(ref string) {
	cmd := exec.Command("git", "checkout", ref)
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		os.Exit(1)
	}
}

func checkIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}

func currentRef() string {
	cmd := exec.Command("git", "log", "--pretty=%H", "-z")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Run()
	return next(&buf)
}

func next(buf *bytes.Buffer) string {
	c, _ := buf.ReadBytes(0)
	c = bytes.TrimSuffix(c, []byte("\000"))
	return string(c)
}

func storeRef() bool {
	refname, err := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD").Output()
	checkIfError(err)
	if strings.Compare(string(refname), "HEAD\n") != 0 {
		err := ioutil.WriteFile(cfgFile, refname, 0666)
		checkIfError(err)
		return true
	}
	return false
}

func loadRef() string {
	dat, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return ""
	}
	dat = bytes.TrimSuffix(dat, []byte("\n"))
	return string(dat)
}
