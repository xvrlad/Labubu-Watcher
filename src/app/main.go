package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
	"github.com/joho/godotenv"
)

func main() {
	// Load env
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal("Error loading .env file:", err)
	}
	url := os.Getenv("POP_MART_URL")

	// Create context with debug and timeout
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // <-- try with headless OFF to see the browser
	)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithDebugf(log.Printf))
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	var productNames []string

	log.Println("Navigating to:", url)

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`div.index_productGrid__zN2jL`),
		chromedp.Evaluate(`(() => {
			const products = Array.from(document.querySelectorAll('.index_productItemContainer__rDwtr'));
			const inStock = [];
			products.forEach(item => {
				const nameEl = item.querySelector('.index_itemTitle__WaT6_');
				const tagEl = item.querySelector('.index_tag__E64FE span');
				const name = nameEl ? nameEl.innerText.trim() : 'Unnamed';
				const outOfStock = tagEl && tagEl.innerText.includes('OUT OF STOCK');
				if (!outOfStock) inStock.push(name);
			});
			return inStock;
		})()`, &productNames),
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n✅ In-stock products:")
	for _, name := range productNames {
		fmt.Println("-", strings.TrimSpace(name))
	}
}
