package music

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/lithdew/youtube"
)

type youtubeDLEntry struct {
	Query        string   `json:"-"`
	ID           string   `json:"id"`
	Title        string   `json:"title"`
	URL          string   `json:"url"`
	ExtractorKey string   `json:"extractor_key"`
	IsLive       bool     `json:"is_live"`
	Duration     *float64 `json:"duration"`
}

type youtubeDLResponse []youtubeDLEntry

func GetYoutubeDLResponse(query string) (youtubeDLResponse, error) {
	resp := youtubeDLResponse{}
	url := strings.HasPrefix(query, "http")
	playlist := strings.HasPrefix(query, "https://www.youtube.com/playlist?") || !url
	spotify := strings.HasPrefix(query, "https://open.spotify.com/playlist/")
	if !spotify {
		if url {
			if playlist {
				ps, err := youtube.LoadPlaylist(strings.Split(query, "list=")[1], 0)
				if err != nil {
					return resp, err
				}
				tracks := make([]string, len(ps.Items))
				for i, track := range ps.Items {
					tracks[i] = fmt.Sprintf("https://www.youtube.com/watch?v=%s", track.ID)
				}
				resp = concurrentResolver(tracks, youtubeURLResolve)
			} else {
				e, err := youtubeURLResolve(query)
				if err != nil {
					return resp, err
				}
				resp = append(resp, e)
			}
		} else {
			e, err := youtubeSearch(query)
			if err != nil {
				return resp, err
			}
			resp = append(resp, e)
		}
	} else {
		tracks, err := GetTracks(query)
		if err != nil {
			return resp, err
		}
		resp = concurrentResolver(tracks, youtubeSearch)
	}
	return resp, nil
}

func concurrentResolver(items []string, function func(str string) (youtubeDLEntry, error)) youtubeDLResponse {
	var resp youtubeDLResponse
	var wg sync.WaitGroup
	_resp := make(youtubeDLResponse, 25)
	for i, t := range items {
		if i == 25 {
			break
		}
		wg.Add(1)
		go func(query string, idx int, wg *sync.WaitGroup) {
			defer wg.Done()
			e, err := function(query)
			if err != nil {
				return
			}
			_resp[idx] = e
		}(t, i, &wg)
	}
	wg.Wait()
	for _, e := range _resp {
		if e.URL != "" {
			resp = append(resp, e)
		}
	}
	return resp
}

func youtubeSearch(query string) (youtubeDLEntry, error) {
	sr, err := youtube.Search(query, 0)
	if err != nil {
		return youtubeDLEntry{}, err
	}
	if len(sr.Items) == 0 {
		return youtubeDLEntry{}, fmt.Errorf("Couldn't find search results.")
	}
	return youtubeURLResolve(fmt.Sprintf("https://www.youtube.com/watch?v=%s", sr.Items[0].ID))
}

func youtubeURLResolve(url string) (youtubeDLEntry, error) {
	var res youtubeDLEntry
	args := []string{
		url,
		"--no-playlist",
		"-f", "worstaudio/worst",
		"-R", "5",
		"-J",
		"--ignore-errors",
		"--no-warnings",
		"--ignore-config",
		"--force-ipv4",
	}
	cmd := exec.Command("youtube-dl", args...)
	data, err := cmd.Output()
	if err != nil && err.Error() != "exit status 1" {
		return res, err
	}
	err = json.Unmarshal(data, &res)
	if err != nil {
		return res, err
	}
	if res.URL == "" {
		return res, fmt.Errorf("No URL")
	}
	res.Query = url
	return res, nil
}
