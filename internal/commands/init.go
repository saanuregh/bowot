package commands

import (
	"bowot/internal/config"

	"github.com/andersfylling/disgord"
	"github.com/auttaja/gommand"
)

var cmds []*gommand.Command

var infoCategory = &gommand.Category{
	Name:        "Information",
	Description: "General commands to retrieve info.",
}

var modCategory = &gommand.Category{
	Name:        "Moderation",
	Description: "Commands made to ease managing and moderating servers.",
}

var utilCategory = &gommand.Category{
	Name:        "Utilities",
	Description: "Simple utilities to do with the bot.",
}

var funCategory = &gommand.Category{
	Name:        "Fun",
	Description: "Have some fun.",
}

var ecoCategory = &gommand.Category{
	Name:        "cOwOins",
	Description: "Invest in some cOwOins.",
}

var ownerCategory = &gommand.Category{
	Name:        "Owner Commands",
	Description: "Accessible to bot owner only.",
	PermissionValidators: []func(ctx *gommand.Context) (string, bool){
		func(ctx *gommand.Context) (string, bool) {
			if ctx.Message.Author.ID == disgord.ParseSnowflakeString(config.C.Bot.Owner) {
				return "", true
			}
			return "Only bot owner can access this command", false
		},
	},
}

func Register(router *gommand.Router) {
	for _, v := range cmds {
		router.SetCommand(v)
	}
}
