package main

import (
	"fmt"
	"log"
	"time"

	"github.com/clauderoy790/ko1-eng-players/internal/scraper"
)

func main() {
	fmt.Println("Generating HTML ...")
	if err := scraper.GenerateHTML(); err != nil {
		log.Fatalf("Failed to generate HTML, err: %s at: %v", err.Error(), time.Now())
	}
	fmt.Println("HTML file generated successfully")
}
