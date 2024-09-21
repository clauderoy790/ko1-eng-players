package scraper

import (
	"encoding/json"
	"fmt"
	"os"
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
	}
}

func removeRecentPlayers(players map[string][]Player) {
	for server := range players {
		toRemove := toMap(players[server])
		recent := toMap(recentPlayers[server])
		newRecent := make([]Player, 0)

		// keep players that aren't in toRemove
		for playerName, player := range recent {
			if _, ok := toRemove[playerName]; !ok {
				newRecent = append(newRecent, player)
			}
		}
		recentPlayers[server] = newRecent
	}
}

// remove players that haven't been seen for more than expiry time
func removeExpiredRecentPlayers(now *time.Time, expiry time.Duration) {
	for server, players := range recentPlayers {
		newRecent := make([]Player, 0)

		// if player has been seen for less then the expiry, keep them
		for _, player := range players {
			lastSeen := utils.StringToTime(player.LastSeen)
			if now.Sub(*lastSeen) <= expiry {
				newRecent = append(newRecent, player)
			}
		}

		recentPlayers[server] = newRecent
	}
}

func updateRecentPlayersLastSeenForDisplay(now *time.Time) {
	for _, players := range recentPlayers {
		for _, p := range players {
			t := utils.StringToTime(p.LastSeen)
			p.LastSeen = utils.DeltaDisplayTime(now, t)
		}
	}
}

func saveRecentPlayers() error {
	return utils.SaveJSON(recentPlayersFile, recentPlayers)
}
