package commands

import (
	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "restart",
		Aliases:     []string{},
		Description: "Restart the bot.",
		Category:    ownerCategory,
		Function:    restart,
	})
}

func restart(ctx *gommand.Context) error {
	_, err := ctx.Reply("Restarting")
	return err
}
