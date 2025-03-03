package stash

import (
	"github.com/spf13/cobra"
)

var StashCmd = &cobra.Command{
  Use: "stash",
  Short: "interact with the stash",
}

func init() {
  StashCmd.AddCommand(StashListCmd)
}
