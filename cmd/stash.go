package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stashCmd = &cobra.Command{
  Use: "stash",
  Short: "interact with the stash",
  RunE: func(cmd *cobra.Command, args []string) error {
    fmt.Println("test")
    return nil
  },
}

func init() {
  rootCmd.AddCommand(stashCmd)
}

