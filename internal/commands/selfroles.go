package commands

import (
	"bowot/internal/db"
	"bowot/internal/embeds"
	"context"
	"fmt"

	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "selfroles",
		Description: "Retrieves all selfroles.",
		Category:    infoCategory,
		Function:    selfroles,
	})
}

func selfroles(ctx *gommand.Context) error {
	if tmp, ok := db.Guilds.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*db.Guild)
		srs := guild.SelfRoles
		if len(srs) < 1 {
			_, err := ctx.Reply("No selfroles available, change selfroles regex using `%settings` or add a new one in server settings.")
			return err
		}
		rs, _ := ctx.Session.GetGuildRoles(context.Background(), ctx.Message.GuildID)
		roleMap := make(map[string]string)
		for _, r := range rs {
			for _, sr := range srs {
				if sr == r.ID.String() {
					roleMap[sr] = r.Name
				}
			}
		}
		desc := "Current selfroles are:\n"
		for id, sr := range roleMap {
			desc += fmt.Sprintf("ID: `%v` Name: `%v`\n", id, sr)
		}
		_, err := ctx.Reply(embeds.Info(
			"Selfroles",
			desc,
			"",
		))
		return err
	}
	return fmt.Errorf("Guild not found.")
}
