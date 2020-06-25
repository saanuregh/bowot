package commands

import (
	"bowot/internal/embeds"
	"bowot/internal/utils"
	"fmt"

	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "dice",
		Aliases:     []string{},
		Description: "Throw a dice.",
		Category:    funCategory,
		Function:    dice,
	})
}

func dice(ctx *gommand.Context) error {
	var dice = []int{1, 2, 3, 4, 5, 6}
	n := dice[utils.GetRandomInt(len(dice))]
	_, err := ctx.Reply(embeds.Info(
		"Dice",
		fmt.Sprintf("You rolled a %v.", n),
		"",
	))
	return err
}
