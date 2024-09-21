package scraper

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const lastOnlinePlayerFile = "../../last-online-players.json"

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
	}

	if err = json.Unmarshal(b, &lastOnline); err != nil {
		return lastOnline, fmt.Errorf("error unmarshalling last online json: %w", err)
	}
	return lastOnline, nil
}

func saveLastOnlinePlayers(lastOnlinePlayers *LastOnlinePlayers) error {
	b, err := json.Marshal(lastOnlinePlayers)
	if err != nil {
		return fmt.Errorf("error marshalling players: %w", err)
	}

	if err := os.WriteFile(lastOnlinePlayerFile, b, 0644); err != nil {
		return fmt.Errorf("error saving to file: %s, err: %w", lastOnlinePlayerFile, err)
	}

	return nil
}

// get the players that were online but are now offline
func getOfflinePlayers(lastOnline *LastOnlinePlayers, currentPlayers map[string][]Player) map[string][]Player {
	nowOffline := make(map[string][]Player)

	for server, lOnline := range lastOnline.Players {
		currentMap := toMap(currentPlayers[server])
		lastMap := toMap(lOnline)

		for name, player := range lastMap {
			if _, ok := currentMap[name]; !ok {
				nowOffline[server] = append(nowOffline[server], player)
			}
		}
	}

	return nowOffline
}
