package commands

import (
	"bowot/internal/db"
	"bowot/internal/embeds"
	"fmt"

	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "wakephrases",
		Aliases:     []string{"wake"},
		Description: "Retrieves all wake phrases.",
		Category:    infoCategory,
		Function:    wakephrase,
	})
}

func wakephrase(ctx *gommand.Context) error {
	if tmp, ok := db.Guilds.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*db.Guild)
		wps := guild.WakePhrases
		if len(wps) < 1 {
			_, err := ctx.Reply("No wake phrases available, add one using `%settings`")
			return err
		}
		wpDesc := "Current wake phrases are:\n"
		for i, w := range wps {
			wpDesc += fmt.Sprintf("%v. Phrase: `%v` Reply: `%v`", i+1, w[0], w[1])
			if len(w) == 3 {
				wpDesc += fmt.Sprintf(" Emoji: %v", w[2])
			}
			wpDesc += "\n"
		}
		_, err := ctx.Reply(embeds.Info(
			"Wake Phrases",
			wpDesc,
			"",
		))
		return err
	}
	return fmt.Errorf("Guild not found.")
}
