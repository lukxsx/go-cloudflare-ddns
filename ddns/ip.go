package ddns

import (
	"io"
	"log"
	"net/http"
)

func GetMyIPAddress() string {
	res, err := http.Get("https://api.ipify.org")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(body)
}
