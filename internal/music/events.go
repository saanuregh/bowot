package music

import (
	"context"

	"github.com/andersfylling/disgord"
)

var events = map[string]interface{}{}

func init() {
	events[disgord.EvtVoiceStateUpdate] = func(session disgord.Session, evt *disgord.VoiceStateUpdate) {
		user, err := session.GetUser(context.Background(), evt.UserID)
		if err != nil {
			session.Logger().Error("error getting user")
		}
		if !user.Bot {
			var guild *GuildPlayer
			if tmp, ok := GuildPlayers.Load(evt.GuildID.String()); ok {
				guild = tmp.(*GuildPlayer)
			} else {
				guild = addGuild(evt.GuildID.String())
			}
			if evt.ChannelID.IsZero() {
				guild.RemoveVoiceState(evt.UserID)
			} else {
				vs, _ := guild.GetVoiceState(evt.UserID)
				if vs == nil {
					guild.AddVoiceState(evt.VoiceState)
				} else {
					guild.UpdateVoiceState(evt.UserID, evt.VoiceState)
				}
			}
		}
	}
}
