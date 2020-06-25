package commands

import (
	"bowot/internal/db"
	"bowot/internal/embeds"
	"fmt"

	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "customs",
		Description: "Retrieves all custom commands.",
		Category:    infoCategory,
		Function:    customs,
	})
}

func customs(ctx *gommand.Context) error {
	if tmp, ok := db.Guilds.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*db.Guild)
		ccmds := guild.CustomCommands
		if len(ccmds) < 1 {
			_, err := ctx.Reply("No custom commands available, add one using `%settings`")
			return err
		}
		cusDesc := "Current custom commands are:\n"
		for i, c := range ccmds {
			cusDesc += fmt.Sprintf("%v. Command: `%v` Reply: `%v`\n", i+1, c[0], c[1])
		}
		_, err := ctx.Reply(embeds.Info(
			"Custom Commands",
			cusDesc,
			"",
		))
		return err
	}
	return fmt.Errorf("Guild not found.")
}
