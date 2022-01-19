package xmux

import (
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestRequest(t *testing.T) {
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
	log.Println(string(b))
	response.Body.Close()
}
