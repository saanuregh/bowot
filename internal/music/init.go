package music

import (
	"github.com/andersfylling/disgord"
	"github.com/auttaja/gommand"
)

var client *disgord.Client

func RegisterEvents(c *disgord.Client) {
	client = c
	for i, v := range events {
		client.On(i, v)
	}
}

func RegisterCommands(router *gommand.Router) {
	for _, v := range cmds {
		router.SetCommand(v)
	}
}
