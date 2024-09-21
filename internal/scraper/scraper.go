package scraper

import (
	"encoding/json"
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
	"github.com/clauderoy790/ko1-eng-players/internal/utils"
)


// sometimes server is displayed differently so we keep track of that
var serverMap = map[string][]string{
	"Otuken":    {"Ötüken", "Otuken"},
	"Ergenekon": {"Ergenekon"},
}

const nbRetries = 5
const retryInSeconds = 10
const clientTimeoutSeconds = 30

// scrapeCurrentPlayers scrapes site data to find current players
func scrapeCurrentPlayers() (map[string][]Player, error) {
	players := make(map[string][]Player)
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

	doc.Find("table.ko1-table tbody tr").Each(func(index int, row *goquery.Selection) {
		// check that server is valid
		serverName := strings.TrimSpace(row.Find("td").Eq(0).Text())
		serverKey := ""
		for key, serverNames := range serverMap {
			if slices.Contains(serverNames, serverName) {
				serverKey = key
				break
			}
		}

		// cant find server
		if serverKey == "" {
			fmt.Println("Cannot find server name in server map: ", serverKey)
			return
		}

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

			players[serverKey] = append(players[serverKey], Player{
				Name:      strings.TrimSpace(playerName),
				Location:  strings.TrimSpace(location),
				NationImg: nationImg,
			})
		}
	})

	return players, nil
}

// GenerateHTML generates the HTML page for all servers
func GenerateHTML() error {
	now := utils.Now()

	lastOnline, err := loadLastOnlinePlayers()
	if err != nil {
		return fmt.Errorf("error loading last online players: %w", err)
	}

	if err = loadRecentPlayers(); err != nil {
		return fmt.Errorf("error loading recent players: %w", err)
	}

	// Scrape currentPlayers for each server
	currentPlayers, err := scrapeCurrentPlayers()
	if err != nil {
		return fmt.Errorf("error loading current players: %w", err)
	}

	// remove any player that is currently online from the recent list
	removeRecentPlayers(currentPlayers)

	// add player that are now offline to the recent list
	nowOffline := getOfflinePlayers(&lastOnline, currentPlayers)

	// set the players last seen date to last update
	setLastSeenDate(lastOnline.UpdateTime, nowOffline)

	// add the offline players as revent
	addRecentPlayers(nowOffline)

	// remove any player that has been offline for more than a month
	removeExpiredRecentPlayers(now, time.Hour*24*30)

	// save the current players as the last online
	lastOnline = LastOnlinePlayers{
		UpdateTime: now,
		Players:    currentPlayers,
	}
	if err = saveLastOnlinePlayers(&lastOnline); err != nil {
		return fmt.Errorf("error saving last online players: %w", err)
	}

	if err = saveRecentPlayers(); err != nil {
		return fmt.Errorf("error saving recent players: %w", err)
	}

	// update recent player last seen to be nicely displayed
	updateRecentPlayersLastSeenForDisplay(now)

	var servers []string
	for k, _ := range serverMap {
		servers = append(servers, k)
	}

	// Prepare the page data
	data := PageData{
		UpdatedAt:     utils.TimeToString(now),
		OnlinePlayers: currentPlayers,
		RecentPlayers: recentPlayers,
		Servers:       servers,
	}

	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println("error marshaling PageData: ", data)
	} else {
		fmt.Printf("Using PageData: \n %s \n", string(b))
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
