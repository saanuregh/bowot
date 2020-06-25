package events

import (
	"bowot/internal/db"
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
)

func init() {
	events[disgord.EvtGuildMemberAdd] = func(session disgord.Session, evt *disgord.GuildMemberAdd) {
		if evt.Member.User.Bot {
			return
		}
		if tmp, ok := db.Guilds.Load(evt.Member.GuildID.String()); ok {
			guild := tmp.(*db.Guild)
			mem := guild.AddMember(evt.Member)
			if mem == nil {
				session.Logger().Error(fmt.Errorf(
					"ADD MEMBER FAILED GUILDID=%s, MEMBERID=%s MSG=failed to add member",
					evt.Member.GuildID.String(),
					evt.Member.User.ID.String(),
				))
				return
			}
			if mem.SelfRole != "" {
				err := session.AddGuildMemberRole(
					context.Background(),
					evt.Member.GuildID,
					evt.Member.User.ID,
					snowflake.ParseSnowflakeString(mem.SelfRole),
				)
				if err != nil {
					session.Logger().Error(fmt.Errorf(
						"ADD MEMBER FAILED GUILDID=%s, MEMBERID=%s MSG=failed to add default selfrole",
						evt.Member.GuildID.String(),
						evt.Member.User.ID.String(),
					))
					return
				}
			}
			session.Logger().Info(fmt.Sprintf(
				"ADD MEMBER SUCCESS GUILDID=%s, MEMBERID=%s",
				evt.Member.GuildID.String(),
				evt.Member.User.ID.String(),
			))
		} else {
			session.Logger().Error(fmt.Errorf(
				"ADD MEMBER FAILED GUILDID=%s, MEMBERID=%s MSG=can't find the guild",
				evt.Member.GuildID,
				evt.Member.User.ID,
			))
		}
	}
}
