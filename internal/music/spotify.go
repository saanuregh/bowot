package music

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

func getStringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str[s:], end)
	if s+e < len(str) && e > 0 {
		result = str[s : s+e]
	}
	return
}

func GetTracks(spotifyURL string) (tracks []string, err error) {
	type Artists struct {
		Name string `json:"name"`
	}
	type SpotifyTrack struct {
		Artists []Artists `json:"artists"`
		Name    string    `json:"name"`
	}
	type Items struct {
		Track SpotifyTrack `json:"track,omitempty"`
	}
	type Tracks struct {
		Items []Items `json:"items"`
	}
	type SpotifyData struct {
		Tracks Tracks `json:"tracks"`
	}
	accessToken, err := getAccessToken(spotifyURL)
	if err != nil {
		return
	}
	foo := strings.Split(spotifyURL, "/playlist/")
	if len(foo) < 2 {
		err = fmt.Errorf("could not get id")
		return
	}
	playlistID := strings.Split(foo[1], "/")[0]
	playlistID = strings.Split(playlistID, "?")[0]
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/playlists/"+playlistID+"?type=track%2Cepisode&market=US", nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:71.0) Gecko/20100101 Firefox/71.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Accept-Language", "en")
	req.Header.Set("Referer", "https://open.spotify.com/")
	req.Header.Set("Authorization", "Bearer "+accessToken)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	var data SpotifyData
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		err = errors.Wrap(err, "could not decode spotify data")
	}
	tracks = make([]string, len(data.Tracks.Items))
	if len(tracks) == 0 {
		err = fmt.Errorf("could not find any tracks")
		return
	}
	for i, track := range data.Tracks.Items {
		tracks[i] = fmt.Sprintf("%s - %s", track.Track.Artists[0].Name, track.Track.Name)
	}
	return
}

func getAccessToken(spotifyURL string) (accessToken string, err error) {
	req, err := http.NewRequest("GET", spotifyURL, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:71.0) Gecko/20100101 Firefox/71.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")
	req.Header.Set("Accept-Language", "en-US,en;q=0.5")
	req.Header.Set("Dnt", "1")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", "sp_ab=%7B%7D; sp_landing=http%3A%2F%2Fopen.spotify.com%2Fplaylist%2F37i9dQZF1EtsXGZhBtSWWl; sp_t=c695ff90921aafb17baa61ea6c01c2f8")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Cache-Control", "max-age=0")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
	accessToken = getStringInBetween(string(bodyBytes), `"accessToken":"`, `"`)
	if len(accessToken) < 3 {
		err = fmt.Errorf("got no access token")
	}
	return
}
