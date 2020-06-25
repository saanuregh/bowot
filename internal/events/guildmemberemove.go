package events

import (
	"bowot/internal/db"
	"fmt"

	"github.com/andersfylling/disgord"
)

func init() {
	events[disgord.EvtGuildMemberRemove] = func(session disgord.Session, evt *disgord.GuildMemberRemove) {
		if evt.User.Bot {
			return
		}
		if tmp, ok := db.Guilds.Load(evt.GuildID.String()); ok {
			guild := tmp.(*db.Guild)
			usr := guild.RemoveMember(evt.User.ID.String())
			if usr == nil {
				session.Logger().Error(fmt.Errorf(
					"REMOVE MEMBER FAILED GUILDID=%s, MEMBERID=%s MSG=failed to delete member",
					evt.GuildID.String(),
					evt.User.ID.String(),
				))
				return
			}
			session.Logger().Info(fmt.Sprintf(
				"REMOVE MEMBER SUCCESS GUILDID=%s, MEMBERID=%s",
				evt.GuildID.String(),
				evt.User.ID.String(),
			))
		} else {
			session.Logger().Error(fmt.Errorf(
				"REMOVE MEMBER FAILED GUILDID=%s, MEMBERID=%s MSG=can't find the guild",
				evt.GuildID.String(),
				evt.User.ID.String(),
			))
		}
	}
}
