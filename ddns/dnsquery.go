package ddns

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type DNSResponse struct {
	Answer []struct {
		Name string `json:"name"`
		Type int    `json:"type"`
		TTL  int    `json:"TTL"`
		Data string `json:"data"`
	} `json:"Answer"`
}

func DNSQuery(domain string) {
	req, err := http.NewRequest("GET", "https://1.1.1.1/dns-query?name="+domain+"&type=A", nil)
	req.Header.Set("accept", "application/dns-json")
	if err != nil {
		log.Fatalln(err)
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var dnsResponse DNSResponse

	err = json.Unmarshal([]byte(body), &dnsResponse)
	if err != nil {
		log.Fatalln(err)
	}

	for _, answer := range dnsResponse.Answer {
		fmt.Println("Name:", answer.Name)
		fmt.Println("Type:", answer.Type)
		fmt.Println("TTL:", answer.TTL)
		fmt.Println("Data:", answer.Data)
	}
}
