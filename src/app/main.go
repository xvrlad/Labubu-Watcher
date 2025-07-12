package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"strings"
)

func main() {
	url := os.Getenv("POP_MART_PRODUCTS_URL");

	// 1. Fetch products.json
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error fetching JSON:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading response body:", err)
	}

	// 2. Parse JSON
	var data Products
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal("Error parsing JSON:", err)
	}

	// 3. Filter for Labubu in stock
	found := false
	var inStockProducts []string

	for _, product := range data.Products {
		if strings.Contains(strings.ToLower(product.Handle), "labubu") {
			for _, variant := range product.Variants {
				if variant.Available {
					found = true
					inStockProducts = append(inStockProducts, fmt.Sprintf("https://www.popmart.nz/products/%s", product.Handle))
					break
				}
			}
		}
	}

	if found {
		// 4. Send email notification
		subject := "Labubu in stock!"
		body := fmt.Sprintf("Labubu restock detected! Links:\n\n%s", strings.Join(inStockProducts, "\n"))

		err := sendEmail(subject, body)
		if err != nil {
			log.Fatal("Failed to send email:", err)
		}
		fmt.Println("Email sent! 🎉")
	} else {
		fmt.Println("No Labubu restock found.")
	}
}

// sendEmail sends a simple email using Gmail SMTP
func sendEmail(subject, body string) error {
	from := os.Getenv("GMAIL_USER")
	password := os.Getenv("GMAIL_APP_PASSWORD")
	to := os.Getenv("GMAIL_USER") // send to yourself

	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	msg := []byte("To: " + to + "\r\n" +
		"Subject: " + subject + "\r\n\r\n" +
		body + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)

	return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, msg)
}