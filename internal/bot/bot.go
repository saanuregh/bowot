package bot

import (
	"bowot/internal/commands"
	"bowot/internal/config"
	"bowot/internal/cron"
	"bowot/internal/db"
	"bowot/internal/embeds"
	"bowot/internal/events"
	"bowot/internal/logger"
	"bowot/internal/music"
	"fmt"
	"io"

	"context"

	"github.com/andersfylling/disgord"
	"github.com/auttaja/gommand"
)

func Start() {
	router := gommand.NewRouter(&gommand.RouterConfig{
		PrefixCheck: gommand.MultiplePrefixCheckers(
			func(ctx *gommand.Context, r io.ReadSeeker) bool {
				guildID := ctx.Message.GuildID.String()
				if tmp, ok := db.Guilds.Load(guildID); ok {
					guild := tmp.(*db.Guild)
					bytes := []byte(guild.Prefix)
					l := len(bytes)
					if res := prefixIterator(l, bytes, r); !res {
						return false
					}
					ctx.Prefix = guild.Prefix
					return true
				}
				return false
			},
		),
		Middleware: []func(ctx *gommand.Context) error{func(ctx *gommand.Context) error {
			ctx.Session.Logger().Info(
				fmt.Sprintf(
					"COMMAND CALLED NAME=%s GUILD=%s USER=%s ARGS=%s",
					ctx.Command.GetName(),
					ctx.Message.GuildID,
					ctx.Message.Author.ID,
					ctx.RawArgs,
				),
			)
			return nil
		}},
	})

	client := disgord.New(disgord.Config{
		BotToken:           config.C.Bot.Token,
		Logger:             logger.CustomLogger,
		ProjectName:        "bowot",
		LoadMembersQuietly: true,
	})

	router.AddErrorHandler(func(ctx *gommand.Context, err error) bool {
		switch err.(type) {
		case *gommand.CommandNotFound, *gommand.CommandBlank:
			return true
		case *gommand.InvalidTransformation:
			ctx.Reply(embeds.Error("Invalid Type", err, false))
			return true
		case *gommand.IncorrectPermissions:
			ctx.Reply(embeds.Error("Missing Permissions", err, false))
			return true
		case *gommand.InvalidArgCount:
			ctx.Reply(embeds.Error("Missing Arguments", err, false))
			return true
		case *gommand.PanicError:
			ctx.Session.Logger().Error(err)
			ctx.Reply(embeds.Error("Panic", err, true))
			return false
		default:
			ctx.Session.Logger().Error(err)
			ctx.Reply(embeds.Error("Handled Error:", err, true))
			return false
		}
	})

	router.CustomCommandsHandler = func(ctx *gommand.Context, cmdname string, r io.ReadSeeker) (bool, error) {
		if tmp, ok := db.Guilds.Load(ctx.Message.GuildID.String()); ok {
			guild := tmp.(*db.Guild)
			customCommands := guild.CustomCommands
			if customCommands == nil {
				return false, nil
			}
			for _, cmd := range customCommands {
				if cmdname == cmd[0] {
					_, err := ctx.Reply(cmd[1])
					return true, err
				}
			}
		}
		return false, nil
	}

	commands.Register(router)
	events.Register(client)

	music.RegisterCommands(router)
	music.RegisterEvents(client)

	cron.Register(client)

	router.Hook(client)

	err := client.StayConnectedUntilInterrupted(context.Background())
	if err != nil {
		panic(err)
	}
}

func prefixIterator(l int, bytes []byte, r io.ReadSeeker) bool {
	i := 0
	ob := make([]byte, 1)
	for i != l {
		_, err := r.Read(ob)
		if err != nil {
			return false
		}
		if ob[0] != bytes[i] {
			return false
		}
		i++
	}
	return true
}
