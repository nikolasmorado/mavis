package cmd

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var errorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("9")).
	Bold(true)

var rootCmd = &cobra.Command{
	Use:           "mavis",
	Short:         "curl that stores your queries",
	Args:          cobra.ExactArgs(2),
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		res, err := mavisFetch(args[1], args[0])
		if err != nil {
			return err
		}

		fmt.Println(res)
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, errorStyle.Render(err.Error()))
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mavis.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func mavisFetch(u, m string) (ret string, err error) {
	if !(m == "get" || m == "post") {
		return "", errors.New("method not supported")
	}

	if u == "" {
		return "", errors.New("url can not be nil")
	}

	mU := strings.ToUpper(m)

	client := &http.Client{}

	req, err := http.NewRequest(mU, u, nil)

	if err != nil {
		return "", errors.New(err.Error())
	}

  res, err := client.Do(req)

	if err != nil {
		return "", errors.New(err.Error())
	}

  defer res.Body.Close()

  body, err := io.ReadAll(res.Body)

	return string(body), nil
}
