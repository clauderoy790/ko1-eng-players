package scraper

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/clauderoy790/ko1-eng-players/internal/utils"
)

const lastOnlinePlayerFile = "last-online-players.json"

type LastOnlinePlayers struct {
	UpdateTime *time.Time
	Players    map[string][]Player
}

func loadLastOnlinePlayers() (LastOnlinePlayers, error) {
	lastOnline := LastOnlinePlayers{
		Players: make(map[string][]Player),
	}
	b, err := os.ReadFile(lastOnlinePlayerFile)
	if err != nil {
		fmt.Printf("error loading last online players: %s\n", err.Error())
		return lastOnline, nil
	}

	if err = json.Unmarshal(b, &lastOnline); err != nil {
		return lastOnline, fmt.Errorf("error unmarshalling last online json: %w", err)
	}

	// set last seen on players
	lastSeen := utils.TimeToString(lastOnline.UpdateTime)
	for k := range lastOnline.Players {
		for i := range lastOnline.Players[k] {
			lastOnline.Players[k][i].LastSeen = lastSeen
		}
	}

	return lastOnline, nil
}

// get the players that were online but are now offline
func getOfflinePlayers(lastOnline *LastOnlinePlayers, currentPlayers map[string][]Player) map[string][]Player {
	nowOffline := make(map[string][]Player)

	for server, lOnline := range lastOnline.Players {
		currentMap := toMap(currentPlayers[server])
		lastMap := toMap(lOnline)

		for name, player := range lastMap {
			if _, ok := currentMap[name]; !ok {
				fmt.Println("player has gone offline: ", player)
				nowOffline[server] = append(nowOffline[server], player)
			}
		}
		nowOffline[server] = removeDuplicates(nowOffline[server])
	}

	return nowOffline
}

func saveLastOnlinePlayers(lastOnlinePlayers *LastOnlinePlayers) error {
	path := os.Getenv("LAST_ONLINE_PLAYERS")
	if path == "" {
		path = lastOnlinePlayerFile
	}
	return utils.SaveJSON(path, lastOnlinePlayers)
}
