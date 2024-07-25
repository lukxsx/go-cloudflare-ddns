package main

import (
	"fmt"
	"log"

	"github.com/lukxsx/go-cloudflare-ddns/ddns"
)

func main() {
	ip, err := ddns.GetMyIPAddress()
	if err != nil {
		log.Fatalln("Failed to query for IP:", err)
	}
	fmt.Println("IP:", ip)
	ddns.DNSQuery("example.com")
}
