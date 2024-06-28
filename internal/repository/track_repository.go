package internal

import (
	"database/sql"
	"fmt"

	"github.com/Darth-Tenebros/Melodic-Visions/internal/model"
)

type TrackRepository struct {
	DB *sql.DB
}

func NewTrackRepository(db *sql.DB) *TrackRepository {
	return &TrackRepository{
		DB: db,
	}
}

func (track *TrackRepository) InsertTrack(tableName string, info model.Track) (bool, error) {

	sqlStatement, err := track.DB.Prepare(fmt.Sprintf(`INSERT INTO %s (track_name, album_name, duration, artist_name) VALUES (?, ?, ?, ?)`, tableName))
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}
	defer sqlStatement.Close()

	_, err = sqlStatement.Exec(info.TrackName, info.AlbumName, info.Duration, info.ArtistName)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}

	return true, nil
}

func (track *TrackRepository) GetTracks(tableName string) ([]model.Track, error) {
	rows, err := track.DB.Query(fmt.Sprintf("SELECT track_name, album_name, duration, artist_name FROM %s", tableName))
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()

	var tracks []model.Track
	for rows.Next() {
		var song model.Track
		err = rows.Scan(&song.TrackName, &song.AlbumName, &song.Duration, &song.ArtistName)
		if err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		tracks = append(tracks, song)
	}

	return tracks, nil
}
