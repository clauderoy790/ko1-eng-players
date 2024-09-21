package scraper

import (
	"time"

	"github.com/clauderoy790/ko1-eng-players/internal/utils"
)

func setLastSeenDate(t *time.Time, players map[string][]Player) {
	for _, pl := range players {
		for i, _ := range pl {
			pl[i].LastSeen = utils.TimeToString(t)
		}

	}
}

// generate a map of players using the name as key
func toMap(players []Player) map[string]Player {
	m := make(map[string]Player)

	for _, p := range players {
		m[p.Name] = p
	}

	return m
}
