package commands

import (
	"bowot/internal/embeds"

	"github.com/auttaja/gommand"
	"github.com/parnurzeal/gorequest"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "chuck",
		Aliases:     []string{},
		Description: "Gets random Chuck Norris joke.",
		Category:    funCategory,
		Function:    chuck,
	})
}

func chuck(ctx *gommand.Context) error {
	var f struct {
		Joke string `json:"value"`
	}
	_, _, errs := gorequest.New().Get("https://api.chucknorris.io/jokes/random").EndStruct(&f)
	if len(errs) > 0 {
		return errs[0]
	}
	_, err := ctx.Reply(embeds.Info(
		"Chuck here",
		f.Joke,
		"",
	))
	return err
}
