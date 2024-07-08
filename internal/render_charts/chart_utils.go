package render_charts

import "github.com/Darth-Tenebros/Melodic-Visions/internal/model"

func AggregateArtistTotalTracks(tracks []model.Track) map[string]int {
	top := make(map[string]int)

	for _, track := range tracks {

		_, ok := top[track.ArtistName]
		if ok {
			top[track.ArtistName] = top[track.ArtistName] + 1
		} else {
			top[track.ArtistName] = 1
		}
	}

	return top
}
