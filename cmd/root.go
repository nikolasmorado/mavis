package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var errorStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("9")).
	Bold(true)

var rootCmd = &cobra.Command{
	Use:           "mavis [method] [url]",
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
var stash bool
var name string
var cookies []string

func init() {
	rootCmd.Flags().StringVarP(&data, "data", "d", "", "Data to pass along with the request")
	rootCmd.Flags().StringSliceVarP(&headers, "header", "H", []string{}, "Headers to be sent with the query formatted as Name:Value ")
	rootCmd.Flags().BoolVar(&stash, "stash", false, "Indicates the request should be stashed")
	rootCmd.Flags().StringVarP(&name, "name", "n", "", "Name to stash the request as, can use slashes to denote directories")
  rootCmd.Flags().StringSliceVarP(&cookies, "cookie", "b", []string{}, "Cookie to send with the request, formatted as Cookie=Value")
}

func validMethod(m string) bool {
	vm := map[string]bool{
		"get":    true,
		"head":   true,
		"post":   true,
		"put":    true,
		"patch":  true,
		"delete": true,
		"trace":  true,
	}

	_, ok := vm[m]

	return ok
}

func fixScheme(u string) string {
	_, _, f := strings.Cut(u, "://")

	if !f {
		return strings.Join([]string{"http://", u}, "")
	}

	return u
}

func getFileName(n string) (*string, *string) {
	ls := strings.LastIndexFunc(n, func(r rune) bool {
		return r == '/'
	})

	if ls == -1 {
		return nil, &n
	}

	dir := n[:ls+1]
	nm := n[ls+1:]

	return &dir, &nm
}

func stashRequest(u, m, n, d string, h, c []string) error {
	md, fname := getFileName(n)

	if fname == nil {
		return errors.New("invalid file name")
	}

	dir := filepath.Join(os.Getenv("HOME"), ".mavis", *md)

	err := os.MkdirAll(dir, 0666)

	if err != nil {
		return err
	}

	hm := make(map[string]string)
	for _, hi := range h {
		parts := strings.SplitN(hi, ":", 2)
		if len(parts) == 2 {
			hm[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	cm := make(map[string]string)
	for _, ci := range c {
		parts := strings.SplitN(ci, "=", 2)
		if len(parts) == 2 {
			cm[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}

	buf := new(bytes.Buffer)

	err = toml.NewEncoder(buf).Encode(map[string]any{
		"request": map[string]any{
			"url":     u,
			"method":  m,
			"data":    d,
			"headers": hm,
			"cookies": cm,
		},
		"settings": map[string]any{
			"fileOutput":      false,
			"outputTransformer": "",
			"outputCookies":   false,
		},
	})

  if err != nil {
    return err
  }

  fp := filepath.Join(dir, *fname+".toml")

	return os.WriteFile(fp, buf.Bytes(), 0666)
}

func mavisFetch(u, m string) (string, error) {
	if ok := validMethod(m); !ok {
		return "", errors.New("method not supported")
	}

	if u == "" {
		return "", errors.New("url can not be nil")
	}

	uF := fixScheme(u)

	re := strings.NewReader(data)

	mU := strings.ToUpper(m)

	client := &http.Client{}

	req, err := http.NewRequest(mU, uF, re)

	if err != nil {
		return "", err
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
		return "", err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return "", err
	}

	if stash {
		err = stashRequest(mU, uF, name, data, headers, []string{})
	}

	return string(body), err
}
