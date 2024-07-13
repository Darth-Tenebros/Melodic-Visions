package model

type Track struct {
	TrackName  string
	AlbumName  string
	Duration   int
	ArtistName string
	Id         string
}

type ArtistDuration struct {
	ArtistName  string
	ArtistCount int
}
