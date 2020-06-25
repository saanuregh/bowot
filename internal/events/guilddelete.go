package events

import (
	"bowot/internal/db"
	"fmt"

	"github.com/andersfylling/disgord"
)

func init() {
	events[disgord.EvtGuildDelete] = func(session disgord.Session, evt *disgord.GuildDelete) {
		guildID := evt.UnavailableGuild.ID
		if tmp, ok := db.Guilds.Load(guildID.String()); ok {
			g := tmp.(*db.Guild)
			g = db.RemoveGuild(guildID.String())
			if g == nil {
				session.Logger().Error(fmt.Errorf("REMOVE GUILD FAILED ID=%s MSG=can't delete", guildID))
				return
			}
			session.Logger().Info(fmt.Sprintf("REMOVE GUILD SUCCESS ID=%s", g.ID))
		} else {
			session.Logger().Error(fmt.Errorf("REMOVE GUILD FAILED ID=%s MSG=doesn't exist", guildID))
		}
	}
}
