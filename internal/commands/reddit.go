package commands

import (
	"bowot/internal/config"
	"bowot/internal/embeds"
	"bowot/internal/utils"
	"fmt"

	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "meme",
		Aliases:     []string{"dank"},
		Description: "Random memes from reddit.",
		Category:    funCategory,
		Function:    meme,
	})
	cmds = append(cmds, &gommand.Command{
		Name:        "copypasta",
		Aliases:     []string{},
		Description: "Random copypasta from reddit.",
		Category:    funCategory,
		Function:    copypasta,
	})
	cmds = append(cmds, &gommand.Command{
		Name:        "whoosh",
		Aliases:     []string{},
		Description: "Random whoosh from reddit.",
		Category:    funCategory,
		Function:    whoosh,
	})
}

func copypasta(ctx *gommand.Context) error {
	copypasta, err := utils.GetRandomPost(config.C.Reddit.CopyPasta, false)
	if err != nil {
		return err
	}
	_, err = ctx.Reply(embeds.Info(
		"Random copypasta",
		copypasta.Text,
		fmt.Sprintf("%v | 🔼: %v 🔽: %v ", copypasta.SubredditNamePrefixed, copypasta.Ups, copypasta.Downs),
	))
	return err
}

func meme(ctx *gommand.Context) error {
	meme, err := utils.GetRandomPost(config.C.Reddit.Meme, true)
	if err != nil {
		return err
	}
	_, err = ctx.Reply(embeds.InfoImage(
		"Random Meme",
		meme.Text,
		fmt.Sprintf("%v | 🔼: %v 🔽: %v ", meme.SubredditNamePrefixed, meme.Ups, meme.Downs),
		meme.URL,
	))
	return err
}

func whoosh(ctx *gommand.Context) error {
	whoosh, err := utils.GetRandomPost(config.C.Reddit.Whoosh, true)
	if err != nil {
		return err
	}
	_, err = ctx.Reply(embeds.InfoImage(
		"Random whoosh",
		whoosh.Text,
		fmt.Sprintf("%v | 🔼: %v 🔽: %v ", whoosh.SubredditNamePrefixed, whoosh.Ups, whoosh.Downs),
		whoosh.URL,
	))
	return err
}
