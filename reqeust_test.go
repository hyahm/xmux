package xmux

import (
	"io"
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/hyahm/golog"
)

func TestRequest(t *testing.T) {
	defer golog.Sync()
	cli := http.Client{
		// Transport: &http.Transport{
		// 	DisableKeepAlives: true,
		// },
	}

	req, err := http.NewRequest("POST", "http://localhost:8888/test/form", strings.NewReader(`{"id": 1}`))
	if err != nil {
		log.Fatal(err)
	}
	response, err := cli.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	b, err := io.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	golog.Info(string(b))
	golog.Info(response.StatusCode)
	response.Body.Close()
}
