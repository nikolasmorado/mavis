package cmd

import (
	// "encoding/json"
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

var headers []string

var data string

func init() {
  rootCmd.Flags().StringVarP(&data, "data", "d", "", "Data to pass along with the request")
  rootCmd.Flags().StringSliceVarP(&headers, "headers", "H", []string{}, "Headers to be sent with the query")
}

func mavisFetch(u, m string) (ret string, err error) {
	if !(m == "get" || m == "post") {
		return "", errors.New("method not supported")
	}

	if u == "" {
		return "", errors.New("url can not be nil")
	}

  re := strings.NewReader(data)

	mU := strings.ToUpper(m)

	client := &http.Client{}

	req, err := http.NewRequest(mU, u, re)

	if err != nil {
		return "", errors.New(err.Error())
	}

  for _, h := range headers {
    b, a, f := strings.Cut(h, ":")

    if !f {
      return "", errors.New("header format invalid")
    }
    
    req.Header.Add(b, a)
  }

  res, err := client.Do(req)

	if err != nil {
		return "", errors.New(err.Error())
	}

  defer res.Body.Close()

  body, err := io.ReadAll(res.Body)

	return string(body), nil
}
