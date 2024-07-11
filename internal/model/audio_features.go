package model

type AudioFeature struct {
	Acousticness float64 `json:"acousticness"`
	Danceability float64 `json:"danceability"`
	Valence      float64 `json:"valence"`
}

type AudioFeatures struct {
	AudioFeatures []AudioFeature `json:"audio_features"`
}
