package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
	git "gopkg.in/src-d/go-git.v4"
)

func init() {
	rootCmd.AddCommand(toCmd)
}

var toCmd = &cobra.Command{
	Use:       "to",
	Short:     "Go to specific version. Accepts start, end and a commit SHA",
	Args:      cobra.ExactValidArgs(1),
	ValidArgs: []string{"start"},
	Run: func(cmd *cobra.Command, args []string) {
		if args[0] == "start" {
			oldestCommit()
		}
	},
}

func currentCommit() string {
	return "HEAD"
}

func oldestCommit() {
	r, _ := git.PlainOpen("./")
	logopts := git.LogOptions{}
	cIter, _ := r.Log(&logopts)
	c0, err := cItmer.Next()
	c := c0
	for {
		if err == io.EOF {
			break
		}
		c = c0
		c0, err = cIter.Next()
	}
	fmt.Println(c)
}
