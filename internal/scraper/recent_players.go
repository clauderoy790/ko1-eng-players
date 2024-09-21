package scraper

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/clauderoy790/ko1-eng-players/internal/utils"
)

const recentPlayersFile = "recent-players.json"

var recentPlayers = make(map[string][]Player)

func loadRecentPlayers() error {
	recentPlayers = make(map[string][]Player)
	b, err := os.ReadFile(recentPlayersFile)
	if err != nil {
		fmt.Printf("error loading recent players: %s\n", err.Error())
		return nil
	}

	if err = json.Unmarshal(b, &recentPlayers); err != nil {
		return fmt.Errorf("error unmarshalling json: %w", err)
	}
	return nil
}

func addRecentPlayers(players map[string][]Player) {
	for server, pl := range players {
		recentPlayers[server] = append(recentPlayers[server], pl...)
		recentPlayers[server] = removeDuplicates(recentPlayers[server])
	}
}

func removeRecentPlayers(players map[string][]Player) {
	for server := range players {
		toRemove := toMap(players[server])
		recent := toMap(recentPlayers[server])
		newRecent := make(map[string]Player) // use map to avoid duplicates

		// keep players that aren't in toRemove
		for playerName, player := range recent {
			if _, ok := toRemove[playerName]; !ok {
				newRecent[player.Name] = player
			}
		}
		recentPlayers[server] = toSlice(newRecent)
	}
}

// remove players that haven't been seen for more than expiry time
func removeExpiredRecentPlayers(now *time.Time, expiry time.Duration) {
	for server, players := range recentPlayers {
		newRecent := make(map[string]Player) // use map to avoid duplicates

		// if player has been seen for less then the expiry, keep them
		for _, player := range players {
			lastSeen := utils.StringToTime(player.LastSeen)
			if now.Sub(*lastSeen) <= expiry {
				newRecent[player.Name] = player
			}
		}

		recentPlayers[server] = toSlice(newRecent)
	}
}

func updateRecentPlayersLastSeenForDisplay(now *time.Time) {
	for _, players := range recentPlayers {
		for i := range players {
			t := utils.StringToTime(players[i].LastSeen)
			players[i].LastSeen = utils.DeltaDisplayTime(now, t)
		}
	}
}

func sortRecentPlayers() error {
	for server, players := range recentPlayers {
		var parseErrors []string
		sort.Slice(players, func(i, j int) bool {
			timeI := utils.StringToTime(players[i].LastSeen)
			timeJ := utils.StringToTime(players[j].LastSeen)
			if timeI == nil || timeJ == nil {
				parseErrors = append(parseErrors, fmt.Sprintf("Error parsing time for player(s) %s and/or %s with time: %s and %s\n", players[i].Name, players[j].Name, players[i].LastSeen, players[j].LastSeen))
				return false // Keep original order for unparseable times
			}
			return timeI.After(*timeJ)
		})
		if len(parseErrors) > 0 {
			return fmt.Errorf("errors parsing times for server %s: %s", server, strings.Join(parseErrors, "; "))
		}
	}
	return nil
}

func saveRecentPlayers() error {
	// make sure there are no duplicte
	for k, players := range recentPlayers {
		recentPlayers[k] = removeDuplicates(players)
	}

	if err := sortRecentPlayers(); err != nil {
		fmt.Printf("error sorting recent players: %s\n", err.Error())
	}
	return utils.SaveJSON(recentPlayersFile, recentPlayers)
}
