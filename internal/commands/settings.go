package commands

import (
	botcron "bowot/internal/cron"
	"bowot/internal/db"
	"bowot/internal/embeds"
	"bowot/internal/utils"
	"context"
	"fmt"
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/snowflake/v4"
	"github.com/auttaja/gommand"
	"github.com/lnquy/cron"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:                 "settings",
		Aliases:              []string{},
		Description:          "Configures the bot.",
		Category:             utilCategory,
		PermissionValidators: []gommand.PermissionValidator{gommand.ADMINISTRATOR},
		Function:             settings,
	})
}

func settings(ctx *gommand.Context) error {
	guildID := ctx.Message.GuildID.String()
	if tmp, ok := db.Guilds.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*db.Guild)
		menu := gommand.NewEmbedMenu(embeds.Info("Settings", "", "Wait for all the reactions to show up before selecting."), ctx)
		menu.NewChildMenu(
			&gommand.ChildMenuOptions{
				Embed: embeds.Info(
					"Server prefix",
					"Enter a new prefix.",
					"Current: "+guild.Prefix,
				),
				Button: &gommand.MenuButton{
					Emoji:       "🇦",
					Name:        "Server prefix",
					Description: "The prefix changes how you execute commands. If it's `|`, you'd use `|ping`.",
				},
				AfterAction: func() {
					res := ctx.WaitForMessage(func(_ disgord.Session, msg *disgord.Message) bool {
						return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
					})
					go ctx.Session.DeleteMessage(context.Background(), ctx.Message.ChannelID, res.ID)
					pre := guild.UpdatePrefix(res.Content)
					if pre == nil {
						ctx.Session.Logger().Error("prefix not set")
						ctx.Reply("Prefix didn't update.")
						return
					}
					ctx.Reply("Prefix updated. You can now use commands like so: `" + *pre + "config`.")
				},
			},
		)
		menu.NewChildMenu(
			&gommand.ChildMenuOptions{
				Embed: embeds.Info(
					"Server selfrole regex",
					"Enter a new selfrole regex.",
					"Current: "+guild.SelfRolesRegex,
				),
				Button: &gommand.MenuButton{
					Emoji:       "🇧",
					Name:        "Server selfrole regex",
					Description: "Self roleregex lets the bot find all the self roles in the server.",
				},
				AfterAction: func() {
					res := ctx.WaitForMessage(func(_ disgord.Session, msg *disgord.Message) bool {
						return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
					})
					go ctx.Session.DeleteMessage(context.Background(), ctx.Message.ChannelID, res.ID)
					reg := guild.UpdateSelfRolesRegex(res.Content)
					if reg == nil {
						ctx.Session.Logger().Error("Selfrole regex invalid")
						ctx.Reply("Selfrole regex invalid.")
						return
					}
					roles, err := ctx.Session.GetGuildRoles(context.Background(), snowflake.ParseSnowflakeString(guildID))
					if err != nil {
						ctx.Session.Logger().Error("Error getting the guild roles.")
						ctx.Reply("Error getting the guild roles.")
						return
					}
					selfroles := guild.UpdateSelfRoles(utils.GetGuildSelfRoles(roles, *reg))
					if selfroles == nil {
						ctx.Session.Logger().Error("Selfroles didn't update.")
						ctx.Reply("Selfroles didn't update.")
						return
					}
					ctx.Reply("Selfrole regex updated.")
				},
			},
		)
		selfRoleMap := make(map[string]string)
		roles, _ := ctx.Session.GetGuildRoles(context.Background(), snowflake.ParseSnowflakeString(guildID))
		for _, r1 := range roles {
			for _, r2 := range guild.SelfRoles {
				if r1.ID.String() == r2 {
					selfRoleMap[r2] = r1.Name
				}
			}
		}
		defSelfRoleDesc := "Change your selfrole to?\nType the any of the role id in the chat.\n`ID` - Name\n"
		for id, s := range selfRoleMap {
			defSelfRoleDesc += fmt.Sprintf("* `%s` - %s\n", id, s)
		}
		menu.NewChildMenu(
			&gommand.ChildMenuOptions{
				Embed: embeds.Info(
					"Default server selfrole",
					defSelfRoleDesc,
					"Current: "+guild.DefaultSelfRole,
				),
				Button: &gommand.MenuButton{
					Emoji:       "🇨",
					Name:        "Default server selfrole",
					Description: "Default selfrole that the bot will assign to new members.",
				},
				AfterAction: func() {
					if len(guild.SelfRoles) < 1 {
						ctx.Reply("There are no selfroles available, either add a new one or update the selfrole regex.")
						return
					}
					res := ctx.WaitForMessage(func(_ disgord.Session, msg *disgord.Message) bool {
						return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
					})
					go ctx.Session.DeleteMessage(context.Background(), ctx.Message.ChannelID, res.ID)
					flag := false
					for _, s := range guild.SelfRoles {
						if s == res.Content {
							flag = true
							break
						}
					}
					if flag {
						pre := guild.UpdateDefaultSelfRole(res.Content)
						if pre == nil {
							ctx.Session.Logger().Error("Default selfrole not set")
							ctx.Reply("Self role regex invalid.")
							return
						}
						ctx.Reply("Default selfrole updated.")
					} else {
						ctx.Reply("Selfrole doesn't exist, choose one from existing list.")
					}
				},
			},
		)
		menu.NewChildMenu(
			&gommand.ChildMenuOptions{
				Embed: embeds.Info(
					"Add a custom command",
					"Enter the the command and reply message in the format `command|reply message`.",
					"Command should a single word.",
				),
				Button: &gommand.MenuButton{
					Emoji:       "🇩",
					Name:        "Add a custom command",
					Description: "Add a new custom command to the server.",
				},
				AfterAction: func() {
					res := ctx.WaitForMessage(func(_ disgord.Session, msg *disgord.Message) bool {
						return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
					})
					go ctx.Session.DeleteMessage(context.Background(), ctx.Message.ChannelID, res.ID)
					cmd := strings.Split(strings.TrimSpace(res.Content), "|")
					if len(cmd) != 2 {
						ctx.Reply("Wrong custom command format. Format is `command|reply message`")
						return
					}
					if len(strings.Split(strings.TrimSpace(cmd[0]), " ")) > 1 {
						ctx.Reply("Wrong custom command format, `command` should be a single word. Format is `command|reply message`")
						return
					}
					flag := false
					for _, c := range ctx.Router.GetAllCommands() {
						if c.GetName() == cmd[0] {
							flag = true
							break
						}
					}
					if flag {
						ctx.Reply("That command is reserved, can't be used as a custom command.")
						return
					}
					if guild.AddCustomCommands(cmd) == nil {
						ctx.Reply("Custom command adding failed.")
						ctx.Session.Logger().Error("Custom command not added.")
						return
					}
					ctx.Reply(fmt.Sprintf("Custom command added. You can now use the commands like so: `%v%v`", guild.Prefix, cmd[0]))
				},
			},
		)
		cusDesc := "Current custom commands are:\n"
		for i, c := range guild.CustomCommands {
			cusDesc += fmt.Sprintf("%v. `%v`\n", i+1, c[0])
		}
		menu.NewChildMenu(
			&gommand.ChildMenuOptions{
				Embed: embeds.Info(
					"Remove a custom command",
					cusDesc,
					"Enter the command to remove it.",
				),
				Button: &gommand.MenuButton{
					Emoji:       "🇪",
					Name:        "Remove a custom command",
					Description: "Remove a custom command from the server.",
				},
				AfterAction: func() {
					res := ctx.WaitForMessage(func(_ disgord.Session, msg *disgord.Message) bool {
						return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
					})
					go ctx.Session.DeleteMessage(context.Background(), ctx.Message.ChannelID, res.ID)
					cmd := strings.TrimSpace(res.Content)
					for _, c := range guild.CustomCommands {
						if cmd == c[0] {
							if guild.RemoveCustomCommands(cmd) == nil {
								ctx.Session.Logger().Error("Custom command not removed.")
								ctx.Reply(fmt.Sprintf("Custom command `%v` not removed.", cmd))
								return
							}
							ctx.Reply(fmt.Sprintf("Custom command `%v` removed.", cmd))
							return
						}
					}
					ctx.Reply(fmt.Sprintf("Custom command `%v` doesn't exist.", cmd))
				},
			},
		)
		menu.NewChildMenu(
			&gommand.ChildMenuOptions{
				Embed: embeds.Info(
					"Add a wake phrase",
					"Enter the the wake phrase, reply message and emoji in the format `phrase|reply message` or `phrase|reply message|emoji`.",
					"Emoji is optional.",
				),
				Button: &gommand.MenuButton{
					Emoji:       "🇫",
					Name:        "Add a wake phrase",
					Description: "Add a new wake phrase to the server.",
				},
				AfterAction: func() {
					res := ctx.WaitForMessage(func(_ disgord.Session, msg *disgord.Message) bool {
						return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
					})
					go ctx.Session.DeleteMessage(context.Background(), ctx.Message.ChannelID, res.ID)
					wp := strings.Split(strings.TrimSpace(res.Content), "|")
					if len(wp) > 3 || len(wp) < 2 {
						ctx.Reply("Wrong wake phrase format. Format is `phrase|reply message` or `phrase|reply message|emoji`")
						return
					}
					if guild.AddWakePhrase(wp) == nil {
						ctx.Reply("Wake phrase not added.")
						ctx.Session.Logger().Error("Wake phrase not added.")
						return
					}
					ctx.Reply(fmt.Sprintf("Wake phrase added. bOwOt will reply `%v` with `%v` from now on.", wp[0], wp[1]))
				},
			},
		)
		wakeDesc := "Current wake phrases are:\n"
		for i, w := range guild.WakePhrases {
			wakeDesc += fmt.Sprintf("%v. `%v`\n", i+1, w[0])
		}
		menu.NewChildMenu(
			&gommand.ChildMenuOptions{
				Embed: embeds.Info(
					"Remove a wake phrase",
					wakeDesc,
					"Enter the wake phrase to remove it.",
				),
				Button: &gommand.MenuButton{
					Emoji:       "🇬",
					Name:        "Remove a wake phrase",
					Description: "Remove a wake phrase from the server.",
				},
				AfterAction: func() {
					res := ctx.WaitForMessage(func(_ disgord.Session, msg *disgord.Message) bool {
						return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
					})
					go ctx.Session.DeleteMessage(context.Background(), ctx.Message.ChannelID, res.ID)
					wp := strings.TrimSpace(res.Content)
					for _, w := range guild.WakePhrases {
						if wp == w[0] {
							if guild.RemoveWakePhrase(wp) == nil {
								ctx.Reply(fmt.Sprintf("Wake phrase `%v` not removed.", wp))
								ctx.Session.Logger().Error("Wake phrase not removed.")
								return
							}
							ctx.Reply(fmt.Sprintf("Wake phrase `%v` removed.", wp))
							return
						}
					}
					ctx.Reply(fmt.Sprintf("Wake phrase `%v` doesn't exist.", wp))
				},
			},
		)
		menu.NewChildMenu(
			&gommand.ChildMenuOptions{
				Embed: embeds.Info(
					"Add a thirsty.",
					"Mention the member.",
					"Mention is must.",
				),
				Button: &gommand.MenuButton{
					Emoji:       "🇭",
					Name:        "Add a thirsty",
					Description: "Add a member who is thirsty and needs some reminding.",
				},
				AfterAction: func() {
					res := ctx.WaitForMessage(func(_ disgord.Session, msg *disgord.Message) bool {
						return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
					})
					go ctx.Session.DeleteMessage(context.Background(), ctx.Message.ChannelID, res.ID)
					if len(res.Mentions) != 1 {
						ctx.Reply("One mention required.")
						return
					}
					err := botcron.AddHydrate(res.Mentions[0].ID.String(), guildID)
					if err != nil {
						ctx.Session.Logger().Error(err)
						ctx.Reply(err)
						return
					}
					mem, _ := ctx.Session.GetMember(context.Background(), ctx.Message.GuildID, res.Mentions[0].ID)
					exprDesc, _ := cron.NewDescriptor()
					desc, _ := exprDesc.ToDescription(botcron.HYDRATEEXPR, cron.Locale_en)
					ctx.Reply(fmt.Sprintf("Bot will remind %v %v.", utils.GetGuildName(mem), strings.ToLower(desc)))
				},
			},
		)
		menu.NewChildMenu(
			&gommand.ChildMenuOptions{
				Embed: embeds.Info(
					"Remove a thirsty.",
					"Mention the member.",
					"Mention is must.",
				),
				Button: &gommand.MenuButton{
					Emoji:       "🇮",
					Name:        "Remove a thirsty",
					Description: "Remove a thirsty member.",
				},
				AfterAction: func() {
					res := ctx.WaitForMessage(func(_ disgord.Session, msg *disgord.Message) bool {
						return msg.Author.ID == ctx.Message.Author.ID && msg.ChannelID == ctx.Message.ChannelID
					})
					go ctx.Session.DeleteMessage(context.Background(), ctx.Message.ChannelID, res.ID)
					if len(res.Mentions) != 1 {
						ctx.Reply("One mention required.")
						return
					}
					err := botcron.RemoveHydrate(res.Mentions[0].ID.String(), guildID)
					if err != nil {
						ctx.Session.Logger().Error(err)
						ctx.Reply(err)
						return
					}
					mem, _ := ctx.Session.GetMember(context.Background(), ctx.Message.GuildID, res.Mentions[0].ID)
					ctx.Reply(fmt.Sprintf("%v won't be reminded again.", utils.GetGuildName(mem)))
				},
			},
		)
		_ = ctx.DisplayEmbedMenu(menu)
		return nil
	}
	return fmt.Errorf("Guild not found.")
}
