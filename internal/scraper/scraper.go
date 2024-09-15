package scraper

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Player struct {
	Name      string
	Location  string
	NationImg string
}

type ServerData struct {
	Server  string
	Players []Player
}

type PageData struct {
	UpdatedAt   string
	ServersData []ServerData
	Servers     []string
}

// sometimes server is displayed differently so we keep track of that
var serverMap = map[string][]string{
	"Otuken":    {"Ötüken", "Otuken"},
	"Ergenekon": {"Ergenekon"},
}

const nbRetries = 5
const retryInSeconds = 10
const clientTimeoutSeconds = 30

// scrapePlayers scrapes player data for the selected server
func scrapePlayers(server string) ([]Player, error) {
	var err error
	url := "https://knightonline1.com/?p=kim_nerede"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("Failed to create request: %s", err)
		return nil, fmt.Errorf("Failed to create request: %s", err)
	}

	client := &http.Client{
		Timeout: clientTimeoutSeconds * time.Second,
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	var resp *http.Response
	retries := 0
	for retries = 0; retries < nbRetries; retries++ {
		resp, err = client.Do(req)
		if err != nil {
			log.Fatalf("failed to fetch the page: %s", err)
		}

		if err == nil && resp.StatusCode == http.StatusOK {
			defer resp.Body.Close()
			break
		}

		if resp.StatusCode != http.StatusOK {
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("failed to read response body, error with code: %d", resp.StatusCode)
			}
			bodyString := string(bodyBytes)
			log.Printf("Error: Failed to fetch the page, status code: %d, message: %s", resp.StatusCode, bodyString)
			time.Sleep(retryInSeconds * time.Second) // wait for retry
		}
	}

	if retries >= nbRetries {
		return nil, fmt.Errorf("failed to fetch the page after %d tries: %s", retries, err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var players []Player
	doc.Find("table.ko1-table tbody tr").Each(func(index int, row *goquery.Selection) {
		// check that server is valid
		serverName := strings.TrimSpace(row.Find("td").Eq(0).Text())
		if slices.Contains(serverMap[server], serverName) {
			// Check if the player speaks English by looking for the en.gif flag
			flagImg, exists := row.Find("img").Attr("src")
			if exists && strings.Contains(flagImg, "en.gif") {
				// Extract the player's name and location
				playerName := row.Find("td").Eq(3).Find("a").Text()
				location := row.Find("td").Eq(1).Text()
				nation, _ := row.Find("td").Eq(6).Find("img").Attr("src")
				nationImg := "./internal/ui/karus.gif"
				if strings.Contains(nation, "elmo") {
					nationImg = "./internal/ui/elmo.gif"
				}

				players = append(players, Player{
					Name:      strings.TrimSpace(playerName),
					Location:  strings.TrimSpace(location),
					NationImg: nationImg,
				})
			}
		}
	})

	return players, nil
}

// GenerateHTML generates the HTML page for all servers
func GenerateHTML() error {
	var servers []string
	var serversData []ServerData

	for server, _ := range serverMap {
		servers = append(servers, server)

		// Scrape players for each server
		players, err := scrapePlayers(server)
		if err != nil {
			return fmt.Errorf("error scraping players for server: %s, err: %w", server, err)
		}
		log.Printf("Found %d players for server: %s", len(players), server)

		serversData = append(serversData, ServerData{
			Server:  server,
			Players: players,
		})
	}

	// Prepare the page data
	location, _ := time.LoadLocation("America/New_York") // Load the EST time zone
	data := PageData{
		UpdatedAt:   time.Now().In(location).Format(time.RFC1123),
		ServersData: serversData,
		Servers:     servers,
	}

	// Parse and execute the template
	tmpl, err := template.ParseFiles("internal/ui/template.html")
	if err != nil {
		return fmt.Errorf("error parsing template: %v", err)
	}

	// Create the output file
	file, err := os.Create("index.html")
	if err != nil {
		return fmt.Errorf("could not create HTML file: %v", err)
	}
	defer file.Close()

	// Execute the template with the data
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("error executing template: %v", err)
	}

	return nil
}
