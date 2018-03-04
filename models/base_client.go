package models

import (
	"net/http"
)

type Authenticator interface {
	BuildAuthURL(string) string
	GetType() string
	newClient(string, *http.Request) (Client, error)
}

type Client interface {
	GetPlaylists() (Playlists, error)
	GetOwnPlaylistInfo(string) (*Playlist, error)
	GetPlaylistInfo(string, string) (*Playlist, error)
	CreatePlaylist(string) (*Playlist, error)
	GetTrackByID(string) (*Track, error)
	GetTracks(*Playlist) (Tracks, error)
	InsertTrack(*Playlist, *Track) (bool, error)
}
