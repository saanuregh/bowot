package commands

import (
	"bowot/internal/embeds"
	"runtime"
	"strings"

	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "about",
		Description: "Retrieves bot information.",
		Category:    infoCategory,
		Function:    about,
	})
}

func about(ctx *gommand.Context) error {
	latency, err := ctx.Session.AvgHeartbeatLatency()
	if err != nil {
		return err
	}
	_, err = ctx.Reply(embeds.Info(
		"About "+strings.Title(ctx.BotUser.Username),
		"This is an instance of **[bowot](https://github.com/fjah/bowot)**, an open-source project.",
		"OwO",
		embeds.Field("language", runtime.Version(), true),
		embeds.Field("os", runtime.GOOS, true),
		embeds.Field("latency", latency.String(), true),
		embeds.Field("built with", "[andersfylling/disgord](https://github.com/andersfylling/disgord)", true),
		embeds.Field("using", "[auttaja/gommand](https://github.com/auttaja/gommand)", true),
	))
	return err
}
