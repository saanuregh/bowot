package commands

import (
	"bowot/internal/embeds"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "poll",
		Aliases:     []string{},
		Description: "Create a poll.",
		Category:    utilCategory,
		Function:    poll,
		ArgTransformers: []gommand.ArgTransformer{
			{
				Optional: false,
				Function: gommand.DurationTransformer,
			},
			{
				Optional:  false,
				Remainder: true,
				Greedy:    true,
				Function:  gommand.StringTransformer,
			},
		},
	})
}

func poll(ctx *gommand.Context) error {
	reactions := []string{
		"🇦",
		"🇧",
		"🇨",
		"🇩",
		"🇪",
		"🇫",
		"🇬",
		"🇭",
		"🇮",
		"🇯",
		"🇰",
		"🇱",
		"🇲",
		"🇳",
		"🇴",
		"🇵",
		"🇶",
		"🇷",
		"🇸",
		"🇹",
	}
	duration := ctx.Args[0].(time.Duration)
	par := strings.Split(ctx.Args[1].(string), "|")
	if len(par) != 2 {
		_, err := ctx.Reply(embeds.Info("Poll",
			"Follow the format `poll duration title|option1,option2,option3`.",
			"",
		))
		return err
	}
	title := par[0]
	options := strings.Split(par[1], ",")
	if len(options) < 2 {
		_, err := ctx.Reply(embeds.Info("Poll", "Need more than two options.", ""))
		return err
	}
	des := "Options are:\n"
	for i := range options {
		des += fmt.Sprintf("%s : %s\n", reactions[i], options[i])
	}
	msg, err := ctx.Reply(embeds.Info(fmt.Sprintf("Poll - %s", title), des, ""))
	if err != nil {
		return err
	}
	for i := range options {
		msg.React(context.Background(), ctx.Session, reactions[i])
	}
	var usersReacted []struct {
		id    disgord.Snowflake
		emoji *disgord.Emoji
	}
	ctx.Session.On(disgord.EvtMessageReactionAdd, func(s disgord.Session, evt *disgord.MessageReactionAdd) {
		if evt.MessageID.String() == msg.ID.String() && evt.UserID.String() != ctx.BotUser.ID.String() {
			flag := false
			for _, r := range reactions {
				if r == evt.PartialEmoji.Name {
					flag = true
				}
			}
			if !flag {
				s.DeleteUserReaction(
					evt.Ctx,
					evt.ChannelID,
					evt.MessageID,
					evt.UserID,
					evt.PartialEmoji.Name,
					disgord.IgnoreCache,
				)
			}
			for i, u := range usersReacted {
				if u.id.String() == evt.UserID.String() {
					s.DeleteUserReaction(
						evt.Ctx,
						evt.ChannelID,
						evt.MessageID,
						evt.UserID,
						u.emoji.Name,
						disgord.IgnoreCache,
					)
					usersReacted = append(usersReacted[:i], usersReacted[i+1:]...)
				}
			}
			usersReacted = append(usersReacted, struct {
				id    disgord.Snowflake
				emoji *disgord.Emoji
			}{evt.UserID, evt.PartialEmoji})
		}
	}, &disgord.Ctrl{Duration: duration})
	time.Sleep(duration)

	msg, err = ctx.Session.GetMessage(context.Background(), msg.ChannelID, msg.ID)
	if err != nil {
		return err
	}
	des = ""
	hv := 0
	ho := []string{}
	vt := 0
	for i, r := range msg.Reactions {
		vt += int(r.Count) - 1
		if int(r.Count)-1 >= hv {
			hv = int(r.Count) - 1
			ho = append(ho, options[i])
		}
		des += fmt.Sprintf("%v - %v - %v\n", r.Emoji.Name, options[i], int(r.Count)-1)
	}
	updatedMsg := ctx.Session.UpdateMessage(
		context.Background(),
		msg.ChannelID,
		msg.ID,
	)
	if vt == 0 {
		updatedMsg.SetEmbed(embeds.Info(fmt.Sprintf("Poll - %s", title), "Nobody voted.", ""))
	} else if len(ho) == 1 {
		des += fmt.Sprintf("\n**%v** won with %v or %v%% votes.", ho[0], hv, (hv/vt)*100)
		updatedMsg.SetEmbed(embeds.Info(fmt.Sprintf("Poll - %s", title), des, ""))
	} else {
		des += "\n"
		for i, h := range ho {
			if i == (len(ho) - 1) {
				des += fmt.Sprintf("and **%v** tied with %v votes.", h, hv)
			} else {
				des += fmt.Sprintf("**%v** ", h)
			}
		}
		updatedMsg.SetEmbed(embeds.Info(fmt.Sprintf("Poll - %s", title), des, ""))
	}
	updatedMsg.Execute()
	err = ctx.Session.DeleteAllReactions(context.Background(), msg.ChannelID, msg.ID)
	return err
}
