package services

import (
	"io"
	"log"
	"net/http"
)

func FetchProducts(url string) ([]byte) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error fetching products:", err)
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
		return nil
	}
	return body
}
