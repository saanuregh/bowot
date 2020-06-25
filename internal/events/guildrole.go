package events

import (
	"bowot/internal/db"
	"bowot/internal/utils"
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
)

func init() {
	updateFunc := func(session disgord.Session, evt interface{}) {
		var guildID disgord.Snowflake
		switch evt.(type) {
		case *disgord.GuildRoleCreate:
			guildID = evt.(*disgord.GuildRoleCreate).GuildID
		case *disgord.GuildRoleDelete:
			guildID = evt.(*disgord.GuildRoleDelete).GuildID
		case *disgord.GuildRoleUpdate:
			guildID = evt.(*disgord.GuildRoleUpdate).GuildID
		}
		roles, err := session.GetGuildRoles(context.Background(), guildID)
		if err != nil {
			session.Logger().Error(fmt.Errorf(
				"UPDATE SELFROLES FAILED GUILDID=%s MSG=can't find the guild roles",
				guildID,
			))
			return
		}
		if tmp, ok := db.Guilds.Load(guildID.String()); ok {
			guild := tmp.(*db.Guild)
			selfroles := utils.GetGuildSelfRoles(roles, guild.SelfRolesRegex)
			s := guild.UpdateSelfRoles(selfroles)
			if s == nil {
				session.Logger().Error(fmt.Errorf(
					"UPDATE SELFROLES FAILED GUILDID=%s MSG=can't find the update selfroles",
					guildID,
				))
				return
			}
			for _, m := range guild.Members {
				if m.SelfRole != "" {
					err = session.AddGuildMemberRole(
						context.Background(),
						snowflake.ParseSnowflakeString(guild.ID),
						disgord.ParseSnowflakeString(m.ID),
						snowflake.ParseSnowflakeString(m.SelfRole),
					)
					if err != nil {
						session.Logger().Error(fmt.Errorf(
							"UPDATE SELFROLES FAILED GUILDID=%s USERID=%s ROLEID=%s MSG=can't update role for the user",
							guildID,
							m.ID,
							m.SelfRole,
						))
						return
					}
				}
			}
			session.Logger().Info(fmt.Sprintf(
				"UPDATE SELFROLES SUCCESS GUILDID=%s",
				guildID,
			))
		} else {
			session.Logger().Error(fmt.Errorf(
				"UPDATE SELFROLES FAILED GUILDID=%s MSG=can't find guild",
				guildID,
			))
		}
	}
	events[disgord.EvtGuildRoleCreate] = func(session disgord.Session, evt *disgord.GuildRoleCreate) { go updateFunc(session, evt) }
	events[disgord.EvtGuildRoleDelete] = func(session disgord.Session, evt *disgord.GuildRoleDelete) { go updateFunc(session, evt) }
	events[disgord.EvtGuildRoleUpdate] = func(session disgord.Session, evt *disgord.GuildRoleUpdate) { go updateFunc(session, evt) }
}
