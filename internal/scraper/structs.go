package scraper

type Player struct {
	Name      string `json:"name"`
	Location  string `json:"location"`
	NationImg string `json:"nationImg"`
	LastSeen  string `json:"lastSeen"`
}

type PageData struct {
	UpdatedAt     string              `json:"updatedAt"`
	OnlinePlayers map[string][]Player `json:"onlinePlayers"`
	RecentPlayers map[string][]Player `json:"recentPlayers"`
	Servers       []string            `json:"servers"`
}
