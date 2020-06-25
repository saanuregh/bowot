package commands

import (
	"bowot/internal/db"
	"bowot/internal/embeds"
	"bowot/internal/utils"
	"context"
	"fmt"
	"sort"

	"github.com/andersfylling/disgord"
	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "leaderboard",
		Aliases:     []string{"leader"},
		Description: "cOwOin leaderboard.",
		Category:    ecoCategory,
		Function:    leaderboard,
	})
}

func leaderboard(ctx *gommand.Context) error {
	guildID := ctx.Message.GuildID
	if tmp, ok := db.Guilds.Load(guildID.String()); ok {
		guild := tmp.(*db.Guild)
		members := guild.Members
		sort.SliceStable(members, func(i, j int) bool {
			return members[i].Coins > members[j].Coins
		})
		desc := ""
		for i, m := range members {
			_m, err := ctx.Session.GetMember(context.Background(), guildID, disgord.ParseSnowflakeString(m.ID))
			if err != nil {
				return err
			}
			desc += fmt.Sprintf("%v. %v - %v cOwOins\n", i+1, utils.GetGuildName(_m), m.Coins)
		}
		_, err := ctx.Reply(embeds.Info(
			"cOwOins - Leaderboard",
			desc,
			"",
		))
		return err
	}
	return fmt.Errorf("Guild not found.")
}
