package stash

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/spf13/cobra"
)

type Request struct {
	URL     string            `toml:"url"`
	Method  string            `toml:"method"`
	Data    string            `toml:"data"`
	Headers map[string]string `toml:"headers"`
	Cookies map[string]string `toml:"cookies"`
}

type Settings struct {
	FileOutput        bool   `toml:"fileOutput"`
	OutputTransformer string `toml:"outputTransformer"`
	OutputCookies     bool   `toml:"outputCookies"`
}

type StashedRequest struct {
	Request  Request  `toml:"request"`
	Settings Settings `toml:"settings"`
}

var StashRunCmd = &cobra.Command{
	Use:   "run [name]",
	Short: "run a stashed command",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
    name := args[0]
		fmt.Println("Running request: ", name)
    res, err := reRun(name)
		if err != nil {
			return err
		}

		fmt.Println(res)
		return nil
	},
}

func getRequestData(n string) (*StashedRequest, error) {
	fname := filepath.Join(os.Getenv("HOME"), ".mavis", "requests", n+".toml")

	f, err := os.ReadFile(fname)

	if err != nil {
		return nil, err
	}

	var req StashedRequest

	err = toml.Unmarshal(f, &req)

	if err != nil {
		return nil, err
	}

	return &req, nil
}

// TODO: Add some validation we are just trusting the user right now lol
func reRun(n string) (string, error) {
	sReq, err := getRequestData(n)

	if err != nil {
		return "", err
	}

	re := strings.NewReader(sReq.Request.Data)

	client := &http.Client{}

	req, err := http.NewRequest(
		sReq.Request.Method,
		sReq.Request.URL,
		re,
	)

	if err != nil {
		return "", err
	}

	for _, h := range sReq.Request.Headers {
    data, ok := sReq.Request.Headers[h]

    if !ok {
      continue
    }

		req.Header.Add(h, data)
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

	return string(body), err
}
