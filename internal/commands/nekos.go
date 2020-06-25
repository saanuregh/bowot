package commands

import (
	"bowot/internal/embeds"
	"bowot/internal/utils"
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/auttaja/gommand"
	"github.com/parnurzeal/gorequest"
)

func init() {
	type command struct {
		Name        string
		Description string
		Text        string
	}
	commands := []command{
		{
			Name:        "baka",
			Description: "Call somebody baka.",
			Text:        "**%s** calls **%s** baka",
		},
		{
			Name:        "cuddle",
			Description: "Cuddle somebody.",
			Text:        "**%s** cuddles **%s**",
		},
		{
			Name:        "hug",
			Description: "Hug somebody.",
			Text:        "**%s** hugs **%s**",
		},
		{
			Name:        "kiss",
			Description: "Kiss somebody.",
			Text:        "**%s** kisses **%s**",
		},
		{
			Name:        "pat",
			Description: "Pat somebody.",
			Text:        "**%s** pats **%s**",
		},
		{
			Name:        "poke",
			Description: "Poke somebody.",
			Text:        "**%s** pokes **%s**",
		},
		{
			Name:        "slap",
			Description: "Slap somebody.",
			Text:        "**%s** slaps **%s**",
		},
		{
			Name:        "smug",
			Description: "Smug somebody.",
			Text:        "**%s** smugs at **%s**",
		},
		{
			Name:        "tickle",
			Description: "Tickle somebody.",
			Text:        "**%s** tickles **%s**",
		},
	}
	for _, c := range commands {
		cmds = append(cmds,
			&gommand.Command{
				Name:        c.Name,
				Aliases:     []string{},
				Description: c.Description,
				Category:    funCategory,
				Function: func(ctx *gommand.Context) error {
					var f struct {
						URL string `json:"url"`
					}
					_, _, errs := gorequest.New().Get("https://nekos.life/api/v2/img/" + c.Name).EndStruct(&f)
					if len(errs) > 0 {
						return errs[0]
					}
					member := ctx.Args[0].(*disgord.Member)
					_, err := ctx.Reply(embeds.InfoImage(
						fmt.Sprintf(c.Text, utils.GetGuildName(ctx.Message.Member), utils.GetGuildName(member)),
						"",
						"",
						f.URL,
					))
					return err
				},
				ArgTransformers: []gommand.ArgTransformer{
					{
						Function: gommand.MemberTransformer,
						Optional: false,
					},
				},
			})
	}
	cmds = append(cmds,
		&gommand.Command{
			Name:        "8ball",
			Aliases:     []string{},
			Description: "8ball, need more info?.",
			Category:    funCategory,
			Function:    eightball,
		},
		&gommand.Command{
			Name:        "fact",
			Aliases:     []string{},
			Description: "Gets random facts.",
			Category:    funCategory,
			Function:    fact,
		},
		&gommand.Command{
			Name:        "why",
			Aliases:     []string{},
			Description: "Why?",
			Category:    funCategory,
			Function:    why,
		})
}

func eightball(ctx *gommand.Context) error {
	var f struct {
		Resposne string `json:"response"`
		URL      string `json:"url"`
	}
	_, _, errs := gorequest.New().Get("https://nekos.life/api/v2/8ball").EndStruct(&f)
	if len(errs) > 0 {
		return errs[0]
	}
	_, err := ctx.Reply(embeds.InfoImage(
		"8ball",
		"",
		f.Resposne,
		f.URL,
	))
	return err
}

func fact(ctx *gommand.Context) error {
	var f struct {
		Fact string `json:"fact"`
	}
	_, _, errs := gorequest.New().Get("https://nekos.life/api/v2/fact").EndStruct(&f)
	if len(errs) > 0 {
		return errs[0]
	}
	_, err := ctx.Reply(embeds.Info(
		"Random Facts",
		f.Fact,
		"",
	))
	return err
}

func why(ctx *gommand.Context) error {
	var f struct {
		Why string `json:"why"`
	}
	_, _, errs := gorequest.New().Get("https://nekos.life/api/v2/why").EndStruct(&f)
	if len(errs) > 0 {
		return errs[0]
	}
	_, err := ctx.Reply(embeds.Info(
		"Why?",
		f.Why,
		"",
	))
	return err
}
