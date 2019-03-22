package cmd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

const (
	cfgFile string = ".git-walk"
)

func init() {
	rootCmd.AddCommand(toCmd)
}

var toCmd = &cobra.Command{
	Use:   "to",
	Short: "Go to specific version. Accepts start, end, and next",
	Long: `
NAME
	git-walk to - Navigate iteratively through git history.

SYNOPSIS
	git-walk to [start|end|next]

DESCRIPTION
	Navigate through git history. Argument can be one of [start, end, next].

	'start': checks out the earliest commit in the current history
	'end': checks out the latest commit in the current directory
	'next': checks out the next more recent commit in the current history

	git-walk, when used with 'start' or 'next' checks out a commit, so git HEAD becomes detached.

	Whenever 'git-walk to start' is run, the current reference is save into a .git-walk file. The contents of this file allow for checking out commits in the future of the target commit.

	Note that, one cannot run 'git-walk to start' for the first time while HEAD is detached, as a non-detached reference needs to be saved.

EXAMPLES
	git-walk to start
	Goes to the first commit, by commit time, in the current history.

	git-walk to next
	Goes to the next commit, chronologically, in the current history.

	git-walk to end
	Goes to the saved reference, which is the one saved in .git-walk the first time 'git-walk to start' is run.
	`,
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"start", "end", "next"},
	Run: func(cmd *cobra.Command, args []string) {
		run(args[0])
	},
}

func run(target string) {
	isref := storeRef()
	curr := currentRef()
	if !isref {
		if ref := loadRef(); ref != "" {
			moveTo(ref)
		}
	}
	if tgt := targetRef(target, curr, isref); tgt != "" {
		moveTo(tgt)
	}
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

func moveTo(ref string) {
	err := exec.Command("git", "checkout", ref).Run()
	checkIfError(err)
}

func targetRef(where string, curr string, isref bool) string {
	switch where {
	case "end":
		return ""
	case "start":
		return startCommit()
	case "next":
		if isref {
			return ""
		}
		return nextCommit(curr)
	default:
		return ""
	}
}

func startCommit() string {
	cmd := exec.Command("git", "log", "--reverse", "--pretty=%H", "-z")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Run()
	return string(buf.Next(40))
}

func nextCommit(curr string) string {
	cmd := exec.Command("git", "log", "--reverse", "--pretty=%H", "-z")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Run()
	c := next(&buf)
	for {
		if c == curr || len(c) == 0 {
			break
		}
		c = next(&buf)
	}
	c = next(&buf)
	return c
}

func next(buf *bytes.Buffer) string {
	c, _ := buf.ReadBytes(0)
	c = bytes.TrimSuffix(c, []byte("\000"))
	return string(c)
}

func currentRef() string {
	cmd := exec.Command("git", "log", "--pretty=%H", "-z")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Run()
	return next(&buf)
}

func checkIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
