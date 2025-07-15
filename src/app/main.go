package main

import (
	"context"
	"fmt"
	"log"
	"os"
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

	var inStockProducts, comingSoonProducts []string
	
	log.Println("Navigating to:", url)
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible(`.index_productItemContainer__rDwtr .index_tag__E64FE span`),
		chromedp.Evaluate(`(() => {
			const products = Array.from(document.querySelectorAll('.index_productItemContainer__rDwtr'));
			const inStock = [];
			const comingSoon = [];
			products.forEach(item => {
				const nameEl = item.querySelector('.index_itemTitle__WaT6_');
				const tagEl = item.querySelector('.index_tag__E64FE span');
				const name = nameEl ? nameEl.innerText.trim() : 'Unnamed';
				const tagText = tagEl ? tagEl.innerText.trim().toUpperCase() : '';
				if (tagText.includes('COMING SOON')) {
					comingSoon.push(name);
				} else if (!tagText.includes('OUT OF STOCK')) {
					inStock.push(name);
				}
			});
			return { inStock, comingSoon };
		})()`, &struct {
			InStock    *[]string `json:"inStock"`
			ComingSoon *[]string `json:"comingSoon"`
		}{&inStockProducts, &comingSoonProducts}),
	)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\n✅ In-stock products:")
	for _, name := range inStockProducts {
		fmt.Println("-", name)
	}

	fmt.Println("\n🕒 Coming soon products:")
	for _, name := range comingSoonProducts {
		fmt.Println("-", name)
	}
}
