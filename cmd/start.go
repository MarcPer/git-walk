package cmd

import (
	"bytes"
	"os/exec"

	"github.com/spf13/cobra"
)

func init() {
	toCmd.AddCommand(startCmd)
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Move back to a specific commit, keeping the current reference.",
	Long: `
NAME
	git-walk to start - Checkout first commit in history or given commit, keeping current reference for navigation.

SYNOPSIS
	git-walk to start [<commit>]

DESCRIPTION
	If <commit> argument is not given, go to first commit in current history line. Otherwise, checkout <commit>. Either way, a commit is checked out, so git HEAD becomes detached.

	Whenever 'git-walk to start' is run, the current reference is save into a .git-walk file. The contents of this file allow for checking out commits in the future of the target commit.

	Note that, one cannot run 'git-walk to start' for the first time while HEAD is detached, as a non-detached reference needs to be saved.

EXAMPLES
	git-walk to start
	Goes to the first commit, by commit time, in the current history.

	git-walk to next <commit>
	Checks out commit <commit>.
	`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			runStart("")
		} else {
			runStart(args[0])
		}
	},
}

func runStart(target string) {
	isref := storeRef()
	if !isref {
		if ref := loadRef(); ref != "" {
			moveTo(ref)
		}
	}

	if target != "" {
		moveTo(target)
		return
	}
	if tgt := startCommit(); tgt != "" {
		moveTo(tgt)
	}
}

func startCommit() string {
	cmd := exec.Command("git", "log", "--reverse", "--pretty=%H", "-z")
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Run()
	return next(&buf)
}
