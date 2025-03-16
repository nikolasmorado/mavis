package stash

import (
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var StashRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "removes a specified request from the stash",
	Args:  cobra.ExactArgs(1),
  RunE: func(cmd *cobra.Command, args []string) error {
    name := args[0]

    err := remove(name)

    if err != nil {
      return err
    }

    return nil
  },
}

func remove(n string) error {
	fname := filepath.Join(os.Getenv("HOME"), ".mavis", "requests", n+".toml")

  err := os.RemoveAll(fname)

  if err != nil {
    return err
  }
  
  return nil

}
