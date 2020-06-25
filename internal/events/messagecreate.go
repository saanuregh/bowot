package events

import (
	"bowot/internal/config"
	"bowot/internal/db"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
	"github.com/parnurzeal/gorequest"
)

var lastWakePhrase = make(map[snowflake.Snowflake]time.Time)

func init() {
	wakePhraseDetection := func(session disgord.Session, evt *disgord.MessageCreate) {
		if evt.Message.Author.Bot {
			return
		}
		if _, ok := lastWakePhrase[evt.Message.GuildID]; !ok {
			lastWakePhrase[evt.Message.GuildID] = time.Time{}
		}
		if time.Now().Sub(lastWakePhrase[evt.Message.GuildID]) < 30*time.Second {
			return
		}
		if tmp, ok := db.Guilds.Load(evt.Message.GuildID.String()); ok {
			guild := tmp.(*db.Guild)
			if !strings.HasPrefix(evt.Message.Content, guild.Prefix) && !evt.Message.Author.Bot {
				for _, w := range guild.WakePhrases {
					if regexp.MustCompile(`(\s+|^)` + strings.ToLower(w[0]) + `(\s+|$)`).MatchString(strings.ToLower(evt.Message.Content)) {
						msg, err := session.CreateMessage(
							context.Background(),
							evt.Message.ChannelID,
							disgord.NewMessageByString(w[1]),
						)
						if err != nil {
							session.Logger().Error(fmt.Errorf(
								"WAKEPHRASE FAILED GUILDID=%s MSGID=%s MSG=failed to create reply message",
								evt.Message.GuildID,
								evt.Message.ID,
							))
							return
						}
						if len(w) > 2 {
							err := msg.React(context.Background(), session, w[2])
							if err != nil {
								session.Logger().Error(fmt.Errorf(
									"WAKEPHRASE FAILED GUILDID=%s MSGID=%s MSG=failed to react to created message",
									evt.Message.GuildID,
									evt.Message.ID,
								))
								return
							}
						}
						lastWakePhrase[evt.Message.GuildID] = time.Now()
						session.Logger().Info(fmt.Sprintf(
							"WAKEPHRASE SUCCESS GUILDID=%s MSGID=%s",
							evt.Message.GuildID,
							evt.Message.ID,
						))
					}
				}
			}
		}
	}

	chatbot := func(session disgord.Session, evt *disgord.MessageCreate) {
		if !config.C.Bot.BotLibre.Enabled {
			return
		}
		if evt.Message.Author.Bot {
			return
		}
		if len(evt.Message.Mentions) < 1 {
			return
		}
		for _, m := range evt.Message.Mentions {
			bot, err := session.GetCurrentUser(context.Background())
			if err != nil {
				session.Logger().Error(fmt.Sprintf(
					"CHATBOT FAILED GUILDID=%s MSGID=%s MSG=failed to get bot user",
					evt.Message.GuildID,
					evt.Message.ID,
				))
				return
			}
			if m.ID == bot.ID {
				if tmp, ok := db.Guilds.Load(evt.Message.GuildID.String()); ok {
					guild := tmp.(*db.Guild)
					if !strings.HasPrefix(evt.Message.Content, guild.Prefix) && !evt.Message.Author.Bot {
						content := regexp.MustCompile(`<@!\d*>`).ReplaceAllString(evt.Message.Content, "")
						m := map[string]interface{}{
							"instance":    config.C.Bot.BotLibre.InstanceID,
							"application": config.C.Bot.BotLibre.ApplicationID,
							"message":     strings.TrimSpace(content),
						}
						mJson, err := json.Marshal(m)
						if err != nil {
							session.Logger().Error(fmt.Sprintf(
								"CHATBOT FAILED GUILDID=%s MSGID=%s MSG=failed to marshall json",
								evt.Message.GuildID,
								evt.Message.ID,
							))
							return
						}
						var response struct {
							Message string `json:"message"`
						}
						_, _, errs := gorequest.New().Post("https://www.botlibre.com/rest/json/chat").
							Send(string(mJson)).
							EndStruct(&response)
						if len(errs) > 0 {
							session.Logger().Error(fmt.Sprintf(
								"CHATBOT FAILED GUILDID=%s MSGID=%s MSG=failed get response from botlibre",
								evt.Message.GuildID,
								evt.Message.ID,
							))
							return
						}
						_, err = session.CreateMessage(
							context.Background(),
							evt.Message.ChannelID,
							disgord.NewMessageByString(response.Message),
						)
						if err != nil {
							session.Logger().Error(fmt.Sprintf(
								"CHATBOT FAILED GUILDID=%s MSGID=%s MSG=failed to send message",
								evt.Message.GuildID,
								evt.Message.ID,
							))
						}
						session.Logger().Info(fmt.Sprintf(
							"CHATBOT SUCCESS GUILDID=%s MSGID=%s",
							evt.Message.GuildID,
							evt.Message.ID,
						))
					}
				} else {
					session.Logger().Error(fmt.Sprintf(
						"CHATBOT FAILED GUILDID=%s MSGID=%s MSG=couldn't find the guild",
						evt.Message.GuildID,
						evt.Message.ID,
					))
				}
			}
		}

	}

	events[disgord.EvtMessageCreate] = func(session disgord.Session, evt *disgord.MessageCreate) {
		go wakePhraseDetection(session, evt)
		go chatbot(session, evt)
	}
}
