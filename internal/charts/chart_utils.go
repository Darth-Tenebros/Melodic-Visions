package charts

import "github.com/Darth-Tenebros/Melodic-Visions/internal/model"

func AggregateArtistTotalDurationListened(tracks []model.Track) map[string]int {
	top := make(map[string]int)

	for _, track := range tracks {

		_, ok := top[track.ArtistName]
		if ok {
			top[track.ArtistName] = top[track.ArtistName] + track.Duration
		} else {
			top[track.ArtistName] = track.Duration
		}
	}

	return top
}
