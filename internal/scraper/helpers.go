package scraper

// generate a map of players using the name as key
func toMap(players []Player) map[string]Player {
	m := make(map[string]Player)

	for _, p := range players {
		m[p.Name] = p
	}

	return m
}
