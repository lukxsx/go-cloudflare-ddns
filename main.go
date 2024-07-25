package main

import (
	"fmt"

	"github.com/lukxsx/go-cloudflare-ddns/ddns"
)

func main() {
	fmt.Println("My IP:", ddns.GetMyIPAddress())
	ddns.DNSQuery("example.com")
}
