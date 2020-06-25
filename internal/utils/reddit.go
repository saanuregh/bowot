package utils

import (
	"fmt"
	"regexp"

	"github.com/parnurzeal/gorequest"
)

type RedditPost struct {
	Title                 string  `json:"title"`
	SubredditNamePrefixed string  `json:"subreddit_name_prefixed"`
	Text                  string  `json:"selftext"`
	Downs                 int     `json:"downs"`
	Ups                   int     `json:"ups"`
	Created               float64 `json:"created"`
	URL                   string  `json:"url"`
	NSFW                  bool    `json:"over_18"`
}

func GetRandomPost(subreddits []string, image bool) (*RedditPost, error) {
	r := regexp.MustCompile(`\.(jpeg|jpg|gif|png)$`)
	var f struct {
		Data struct {
			Dist     int `json:"dist"`
			Children []struct {
				Data RedditPost `json:"data,omitempty"`
			} `json:"children"`
		} `json:"data"`
	}
	_, _, errs := gorequest.
		New().
		Get("https://www.reddit.com/r/"+subreddits[GetRandomInt(len(subreddits))]+"/hot/.json").
		AppendHeader("User-Agent", "bowot").
		EndStruct(&f)
	if len(errs) > 0 {
		return nil, errs[0]
	}
	for j := 0; j < f.Data.Dist; j++ {
		i := GetRandomInt(f.Data.Dist)
		if f.Data.Children[i].Data.NSFW {
			continue
		}
		if image {
			if r.FindString(f.Data.Children[i].Data.URL) != "" {
				return &f.Data.Children[i].Data, nil
			}
		} else {
			return &f.Data.Children[i].Data, nil
		}
	}
	return nil, fmt.Errorf("no memes?")
}
