package commands

import (
	"bowot/internal/db"
	"bowot/internal/embeds"
	"context"
	"fmt"
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "user",
		Aliases:     []string{},
		Description: "Configures user settings.",
		Category:    utilCategory,
		Function:    user,
	})
}

func user(ctx *gommand.Context) error {
	guildID := ctx.Message.GuildID.String()
	userId := ctx.Message.Author.ID.String()
	if tmp, ok := db.Guilds.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*db.Guild)
		currentSelfRole := guild.GetSelfRole(userId)
		selfRoleMap := make(map[string]string)
		roles, err := ctx.Session.GetGuildRoles(context.Background(), snowflake.ParseSnowflakeString(guildID))
		if err != nil {
			return err
		}
		for _, r1 := range roles {
			for _, r2 := range guild.SelfRoles {
				if r1.ID.String() == r2 {
					selfRoleMap[r2] = r1.Name
				}
			}
		}
		selfRoleDesc := "Change your selfrole to?\nType the any of the role id below in the chat.\n`ID` - Name\n"
		for id, s := range selfRoleMap {
			selfRoleDesc += fmt.Sprintf("* `%s` - %s\n", id, s)
		}
		menu := gommand.NewEmbedMenu(embeds.Info("User Settings", "", ""), ctx)
		_ = menu.NewChildMenu(
			&gommand.ChildMenuOptions{
				Embed: embeds.Info(
					"Change selfrole",
					selfRoleDesc,
					"Current: "+selfRoleMap[*currentSelfRole],
				),
				Button: &gommand.MenuButton{
					Emoji:       "🌈",
					Name:        "Change selfrole",
					Description: "Change your selfrole OwO.",
				},
				AfterAction: func() {
					res := ctx.WaitForMessage(func(_ disgord.Session, msg *disgord.Message) bool {
						return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
					})
					go ctx.Session.DeleteMessage(context.Background(), ctx.Message.ChannelID, res.ID)
					flag := false
					key := ""
					for k := range selfRoleMap {
						if strings.TrimSpace(res.Content) == k {
							flag = true
							key = k
							break
						}
					}
					if flag {
						ctx.Session.RemoveGuildMemberRole(
							context.Background(),
							snowflake.ParseSnowflakeString(guildID),
							snowflake.ParseSnowflakeString(userId),
							disgord.ParseSnowflakeString(*currentSelfRole),
						)
						guild.UpdateSelfRole(key, userId)
						ctx.Session.AddGuildMemberRole(
							context.Background(),
							snowflake.ParseSnowflakeString(guildID),
							snowflake.ParseSnowflakeString(userId),
							disgord.ParseSnowflakeString(key),
						)
						ctx.Reply(fmt.Sprintf("Selfrole changed to %v.", selfRoleMap[key]))
					} else {
						ctx.Reply("Wrong selfrole selected, try again later.")
					}
				},
			},
		)
		_ = ctx.DisplayEmbedMenu(menu)
		return nil
	}
	return fmt.Errorf("Guild not found.")
}
