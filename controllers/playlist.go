package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/cschen13/spotitube/models"
	"github.com/cschen13/spotitube/utils"
	"github.com/gorilla/mux"

	// "log"
	"net/http"
)

const (
	PAGE_PARAM   = "page"
	CLIENT_PARAM = "client"
)

type PlaylistController struct {
	sessionManager *utils.SessionManager
	// currentUser    *utils.CurrentUserManager
}

func NewPlaylistController(sessionManager *utils.SessionManager) *PlaylistController {
	return &PlaylistController{sessionManager: sessionManager}
}

func (ctrl *PlaylistController) Register(router *mux.Router) {
	router.Handle("/playlists", utils.Handler(ctrl.getPlaylists))
	router.Handle("/playlists/{"+OWNER_ID_PARAM+"}/{"+PLAYLIST_ID_PARAM+"}", utils.Handler(ctrl.getPlaylistInfo))
}

type playlistsClient interface {
	GetPlaylists() (models.Playlists, error)
}

func (ctrl *PlaylistController) getPlaylists(w http.ResponseWriter, r *http.Request) error {
	// user := ctrl.currentUser.Get(r)
	// if user == nil {
	// 	return utils.StatusError{http.StatusUnauthorized, errors.New("getPlaylists: user not logged in")}
	// }

	clientParam := r.URL.Query().Get(CLIENT_PARAM)
	if clientParam == "" {
		clientParam = models.SPOTIFY_SERVICE
	}

	// c := user.GetClient(clientParam)
	c := ctrl.sessionManager.GetClient(r, clientParam)
	if c == nil {
		return utils.StatusError{http.StatusUnauthorized, errors.New(fmt.Sprintf("getPlaylists: no %s client found", clientParam))}
	}

	client, ok := c.(playlistsClient)
	log.Printf("getPlaylists: %v", client)
	if !ok {
		return utils.StatusError{http.StatusMethodNotAllowed, errors.New(fmt.Sprintf("getPlaylists: %s client does not satisfy interface", clientParam))}
	}

	playlists, err := client.GetPlaylists()
	if err != nil {
		return utils.StatusError{http.StatusInternalServerError, err}
	}

	js, err := json.Marshal(playlists)
	if err != nil {
		return utils.StatusError{http.StatusInternalServerError, err}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}

type playlistInfoClient interface {
	GetPlaylistInfo(string, string) (*models.Playlist, error)
}

func (ctrl *PlaylistController) getPlaylistInfo(w http.ResponseWriter, r *http.Request) error {
	// user := ctrl.currentUser.Get(r)
	// if user == nil {
	// 	return utils.StatusError{http.StatusUnauthorized, errors.New("getPlaylists: user not logged in")}
	// }

	clientParam := r.URL.Query().Get(CLIENT_PARAM)
	if clientParam == "" {
		clientParam = models.SPOTIFY_SERVICE
	}

	// c := user.GetClient(clientParam)
	c := ctrl.sessionManager.GetClient(r, clientParam)
	if c == nil {
		return utils.StatusError{http.StatusUnauthorized, errors.New(fmt.Sprintf("getPlaylists: no %s client found", clientParam))}
	}

	client, ok := c.(playlistInfoClient)
	if !ok {
		return utils.StatusError{http.StatusMethodNotAllowed, errors.New(fmt.Sprintf("getPlaylists: %s client cannot be used to get playlist info"))}
	}

	vars := mux.Vars(r)
	ownerId, present := vars[OWNER_ID_PARAM]
	if !present {
		return utils.StatusError{http.StatusUnprocessableEntity, errors.New("getPlaylists: URL is missing an owner ID")}
	}

	playlistId, present := vars[PLAYLIST_ID_PARAM]
	if !present {
		return utils.StatusError{http.StatusUnprocessableEntity, errors.New("getPlaylists: URL is missing a playlist ID")}
	}

	playlist, err := client.GetPlaylistInfo(ownerId, playlistId)
	if err != nil {
		return utils.StatusError{http.StatusInternalServerError, err}
	}

	js, err := json.Marshal(playlist)
	if err != nil {
		return utils.StatusError{http.StatusInternalServerError, err}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return nil
}
