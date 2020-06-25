package commands

import (
	"bowot/internal/db"
	"bowot/internal/embeds"
	"fmt"
	"time"

	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "daily",
		Description: "Get your daily cOwOins.).",
		Category:    ecoCategory,
		Function:    daily,
	})
}

func daily(ctx *gommand.Context) error {
	userId := ctx.Message.Author.ID.String()
	if tmp, ok := db.Guilds.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*db.Guild)
		last := guild.GetLastDaily(userId)
		if last == nil {
			return fmt.Errorf("Can't find the user.")
		}
		if time.Now().Sub(time.Unix(0, *last)) < time.Hour*24 {
			_, err := ctx.Reply(embeds.Info(
				"cOwOins - Daily",
				"You have wait 24 hours to redeem your next daily cOwOins.",
				"",
			))
			return err
		}
		coins := guild.UpdateCoins(100, "+", userId)
		if coins == nil {
			return fmt.Errorf("Can't find the user.")
		}
		last = guild.UpdateLastDaily(userId)
		if last == nil {
			return fmt.Errorf("Can't find the user.")
		}
		_, err := ctx.Reply(embeds.Info(
			"cOwOins - Daily",
			"You have redeemed your daily 100 cOwOins, don't forget to grab tommorow's.",
			fmt.Sprintf("You have %v cOwOins currently.", *coins),
		))
		return err
	}
	return fmt.Errorf("Guild not found.")

}
