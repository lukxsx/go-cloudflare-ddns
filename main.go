package main

import (
	"github.com/lukxsx/go-cloudflare-ddns/ddns"
)

func main() {
	client := &ddns.Client{}
	err := client.Configure()
	if err != nil {
		panic(err)
	}
}
