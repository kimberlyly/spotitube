package models

import (
	"encoding/json"
	"errors"
	youtube "google.golang.org/api/youtube/v3"
	"log"
	"strings"
)

type youtubePlaylist struct {
	obj *youtube.Playlist
}

func (playlist *youtubePlaylist) GetID() string {
	return playlist.obj.Id
}

func (playlist *youtubePlaylist) GetOwnerID() string {
	return playlist.obj.Snippet.ChannelId
}

func (playlist *youtubePlaylist) GetName() string {
	return playlist.obj.Snippet.Title
}

func (playlist *youtubePlaylist) GetURL() string {
	return "https://www.youtube.com/playlist?list=" + playlist.GetID()
}

func (playlist *youtubePlaylist) GetCoverURL() string {
	if thumbnails := playlist.obj.Snippet.Thumbnails; thumbnails != nil {
		return thumbnails.Default.Url
	}
	return ""
}

func (client *youtubeClient) GetPlaylists() (playlists Playlists, err error) {
	return nil, errors.New("youtube: GetPlaylists is unimplemented.")
}

func (client *youtubeClient) GetOwnPlaylistInfo(playlistId string) (*Playlist, error) {
	call := client.Channels.List("id")
	call.Mine(true)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	return client.GetPlaylistInfo(response.Items[0].Id, playlistId)
}

func (client *youtubeClient) GetPlaylistInfo(channelId, playlistId string) (*Playlist, error) {
	call := client.Playlists.List("id,snippet").Id(playlistId)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	if len(response.Items) == 0 {
		return nil, nil
	}

	return NewPlaylist(&youtubePlaylist{response.Items[0]}), nil
}

func (client *youtubeClient) CreatePlaylist(name string) (*Playlist, error) {
	playlist := &youtube.Playlist{}
	properties := map[string]string{
		"snippet.title": name,
	}
	res := createResource(properties)
	if err := json.NewDecoder(strings.NewReader(res)).Decode(&playlist); err != nil {
		log.Printf("youtube: Failed to decode JSON into playlist resource")
		return nil, err
	}

	call := client.Playlists.Insert("id,snippet", playlist)
	response, err := call.Do()
	if err != nil {
		return nil, err
	}

	return NewPlaylist(&youtubePlaylist{response}), nil
}
