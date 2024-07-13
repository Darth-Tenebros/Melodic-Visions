package model

func convertItemToTrack(item Item) Track {
	artistName := ""
	if len(item.Artists) > 0 {
		artistName = item.Artists[0].Name
	}

	return Track{
		TrackName:  item.Name,
		AlbumName:  item.Album.Name,
		Duration:   item.DurationMs,
		ArtistName: artistName,
		Id:         item.Id,
	}
}

func ConvertItemsToTracks(items []Item) ([]Track, []string) {
	tracks := make([]Track, len(items))
	trackIds := make([]string, len(items))
	for i, item := range items {
		tracks[i] = convertItemToTrack(item)
		trackIds[i] = tracks[i].Id
	}
	return tracks, trackIds
}

func RetrieveAudioFeatures(items []AudioFeature) ([]float64, []float64, []float64) {
	var acousticness []float64
	var danceability []float64
	var valence []float64

	for _, feature := range items {
		acousticness = append(acousticness, feature.Acousticness)
		danceability = append(danceability, feature.Danceability)
		valence = append(valence, feature.Valence)
	}

	return acousticness, danceability, valence
}
