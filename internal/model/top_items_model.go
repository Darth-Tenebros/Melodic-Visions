package model

type SpotifyResponse struct {
	Next  string `json:"next"`
	Items []Item `json:"items"`
}

type Item struct {
	Album struct {
		Name string `json:"name"`
	} `json:"album"`
	Name       string   `json:"name"`
	Artists    []Artist `json:"artists"`
	DurationMs int      `json:"duration_ms"`
}

type Artist struct {
	Name string `json:"name"`
}
