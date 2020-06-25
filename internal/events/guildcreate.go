package events

import (
	"bowot/internal/db"
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
)

func init() {
	events[disgord.EvtGuildCreate] = func(session disgord.Session, evt *disgord.GuildCreate) {
		if _, ok := db.Guilds.Load(evt.Guild.ID.String()); !ok {
			guild, err := session.GetGuild(context.Background(), evt.Guild.ID)
			if err != nil {
				session.Logger().Error(fmt.Errorf("ADD GUILD FAILED ID=%s MSG=can't get guilds from discord", evt.Guild.ID))
			}
			members, err := session.GetMembers(
				context.Background(),
				guild.ID,
				&disgord.GetMembersParams{},
			)
			if err != nil {
				session.Logger().Error(fmt.Errorf("ADD GUILD FAILED ID=%s MSG=can't get guild members from discord", evt.Guild.ID))
				return
			}
			g := db.AddGuild(guild, members)
			session.Logger().Info(fmt.Sprintf(
				"ADD GUILD SUCCESS ID=%s NAME=%s N_MEMBERS=%d",
				g.ID,
				guild.Name,
				len(g.Members),
			))
		}
	}
}
