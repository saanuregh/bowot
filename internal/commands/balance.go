package commands

import (
	"bowot/internal/db"
	"bowot/internal/embeds"
	"fmt"

	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "balance",
		Description: "Retrieves your curret cOwOin balance.",
		Category:    ecoCategory,
		Function:    balance,
	})
}

func balance(ctx *gommand.Context) error {
	if tmp, ok := db.Guilds.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*db.Guild)
		coins := guild.GetCoins(ctx.Message.Author.ID.String())
		if coins == nil {
			return fmt.Errorf("Can't find the user.")
		}
		_, err := ctx.Reply(embeds.Info(
			"cOwOins - Balance",
			fmt.Sprintf("You have **%v** cOwOins.", *coins),
			"",
		))
		return err
	}
	return fmt.Errorf("Guild not found.")
}
