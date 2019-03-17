package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/spf13/cobra"
	git "gopkg.in/src-d/go-git.v4"
)

// File cannot be created on the root directory, otherwise it is
// removed upon call to Checkout. See this issue on src-d/go-git:
// https://github.com/src-d/go-git/issues/1026
const (
	cfgFile string = "../.git-walk"
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

	Whenever 'git-walk to start' is run, the current reference is save into the ../.git-walk file. The contents of this file allow for checking out commits in the future of the target commit.

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
	r, err := git.PlainOpen("./")
	checkIfError(err)
	c := currentCommit(r)
	storeRef(r)
	moveToRef(r)
	tgt := targetCommit(r, target, c)
	if tgt == nil {
		return
	}
	moveTo(r, tgt)
}

func currentCommit(r *git.Repository) *object.Commit {
	logopts := git.LogOptions{Order: git.LogOrderCommitterTime}
	cIter, err := r.Log(&logopts)
	checkIfError(err)
	c, err := cIter.Next()
	checkIfError(err)
	return c
}

func storeRef(r *git.Repository) {
	ref, _ := r.Head()
	refname := ref.Name()
	if refname != plumbing.HEAD {
		err := ioutil.WriteFile(cfgFile, []byte(refname), 0666)
		checkIfError(err)
	}
}

func moveToRef(r *git.Repository) {
	dat, err := ioutil.ReadFile(cfgFile)
	checkIfError(err)
	w, err := r.Worktree()
	checkIfError(err)
	checkopts := git.CheckoutOptions{Branch: plumbing.ReferenceName(string(dat)), Force: true}
	err = w.Checkout(&checkopts)
	checkIfError(err)
}

func targetCommit(r *git.Repository, where string, current *object.Commit) *object.Commit {
	var c *object.Commit
	switch where {
	case "start":
		c = oldestCommit(r, nil)
	case "next":
		head, _ := r.Head()
		if current.Hash.String() == head.Hash().String() {
			return nil
		}
		c = oldestCommit(r, current)
	default:
		c = nil
	}
	return c
}

func moveTo(r *git.Repository, c *object.Commit) {
	w, err := r.Worktree()
	checkIfError(err)
	checkopts := git.CheckoutOptions{Hash: c.Hash, Force: true}
	err = w.Checkout(&checkopts)
	checkIfError(err)
}

func oldestCommit(r *git.Repository, current *object.Commit) *object.Commit {
	var c *object.Commit
	if current != nil {
		c = current
	}
	logopts := git.LogOptions{Order: git.LogOrderCommitterTime}
	cIter, err := r.Log(&logopts)
	checkIfError(err)
	c0, err := cIter.Next()
	for {
		if err == io.EOF {
			break
		}
		if current != nil && current.Hash.String() == c0.Hash.String() {
			break
		}
		c = c0
		c0, err = cIter.Next()
	}
	return c
}

func checkIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}
