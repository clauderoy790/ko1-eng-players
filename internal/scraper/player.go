package scraper

import (
	"sort"
	"time"

	"github.com/clauderoy790/ko1-eng-players/internal/utils"
)

func setLastSeenDate(t *time.Time, players map[string][]Player) {
	for _, pl := range players {
		for i := range pl {
			pl[i].LastSeen = utils.TimeToString(t)
		}

	}
}

// remove any duplicate player Name and order's them by player's location/player name
func removeDuplicates(players []Player) []Player {
	m := toMap(players)
	return toSlice(m)
}

// generate a map of players using the name as key
func toMap(players []Player) map[string]Player {
	m := make(map[string]Player)

	for _, p := range players {
		m[p.Name] = p
	}

	return m
}

// returns a slice that's ordered by players location/name
func toSlice(players map[string]Player) []Player {
	sl := make([]Player, 0)

	for _, p := range players {
		sl = append(sl, p)
	}

	sort.Slice(sl, func(i, j int) bool {
		if sl[i].Location != sl[j].Location {
			return sl[i].Location < sl[j].Location
		}
		return sl[i].Name < sl[j].Name
	})

	return sl
}
