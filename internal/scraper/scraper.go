package scraper

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Player struct {
	Name     string
	Location string
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

// scrapePlayers scrapes player data for the selected server
func scrapePlayers(server string) ([]Player, error) {
	url := "https://knightonline1.com/?p=kim_nerede"
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch the page: %v", err)
	}
	defer resp.Body.Close()

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
				players = append(players, Player{
					Name:     strings.TrimSpace(playerName),
					Location: strings.TrimSpace(location),
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

		serversData = append(serversData, ServerData{
			Server:  server,
			Players: players,
		})
	}

	// Prepare the page data
	data := PageData{
		UpdatedAt:   time.Now().Format(time.RFC1123),
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
