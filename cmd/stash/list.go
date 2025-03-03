package stash

import (
	// "fmt"

	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var StashListCmd = &cobra.Command{
	Use:   "list",
	Short: "list out all stashed requests",
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := filepath.Join(os.Getenv("HOME"), ".mavis", "requests")
		fmt.Println("Stashed Requests")
		list(dir, 0, false)
		return nil
	},
}

func list(p string, d int, pl bool) {

	fs, err := os.ReadDir(p)

	if err != nil {
		return
	}

	for i, f := range fs {
		last := i == len(fs)-1
		prefix := ""

		if d > 0 {
			if pl {
        prefix += "    "
      } else {
				prefix += "│   "
			}
		}

    for i := 0; i < 4*(d-1); i++ {
      prefix += " "
    }

    if last {
      prefix += "└──"
    } else {
      prefix += "├──"
    }


		fn := strings.TrimSuffix(f.Name(), ".toml")

    fmt.Println(prefix, fn)

		if !f.IsDir() {
      continue
		}

		np := filepath.Join(p, fn)

		list(np, d+1, last)
	}
}
