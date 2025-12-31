package main

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/Thordata/thordata-go-sdk/examples/internal/example"
)

func main() {
	example.LoadEnv()

	// This example uses Proxy Network. Usually you need THORDATA_PROXY_HOST/PORT (and sometimes username/password).
	if example.SkipIfMissing("THORDATA_PROXY_HOST", "THORDATA_PROXY_PORT") {
		return
	}

	client, err := example.NewClientAllowDummyScraper(60 * time.Second)
	if err != nil {
		panic(err)
	}

	resp, err := client.ProxyGet(context.Background(), "https://httpbin.org/ip", nil)
	if err != nil {
		panic(err)
	}
	defer func() { _ = resp.Body.Close() }()

	b, _ := io.ReadAll(resp.Body)
	fmt.Println("status:", resp.StatusCode)
	fmt.Println(string(b))
}
