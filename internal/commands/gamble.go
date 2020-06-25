package commands

import (
	"bowot/internal/db"
	"bowot/internal/embeds"
	"bowot/internal/utils"
	"fmt"
	"strconv"
	"strings"

	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:            "gamble",
		Description:     "Gamble your cOwOins.).",
		Category:        ecoCategory,
		Function:        gamble,
		ArgTransformers: []gommand.ArgTransformer{{Optional: false, Function: gommand.StringTransformer}},
	})
}

func gamble(ctx *gommand.Context) error {
	GambleMap := []struct {
		val int
		mul int
	}{
		{36, 1},
		{75, 2},
		{90, 3},
		{96, 4},
		{98, 5},
	}
	var err error
	if tmp, ok := db.Guilds.Load(ctx.Message.GuildID.String()); ok {
		guild := tmp.(*db.Guild)
		userId := ctx.Message.Author.ID.String()
		amount := 0
		balance := guild.GetCoins(userId)
		if balance == nil {
			return fmt.Errorf("Can't find the user.")
		}
		amountStr := strings.TrimSpace(ctx.Args[0].(string))
		if amountStr == "all" {
			amount = int(*balance)
		} else {
			amount, err = strconv.Atoi(amountStr)
			if err != nil {
				return fmt.Errorf("Give proper amount value.")
			}
			if amount > int(*balance) {
				_, err := ctx.Reply(embeds.Info(
					"cOwOins - Gamble",
					"You don't have that much balance cOwOins.",
					fmt.Sprintf("You have %v cOwOins currently.", *balance),
				))
				return err
			}
		}
		randVal := utils.GetRandomIntDist(100, 2)
		mul := 0
		for _, g := range GambleMap {
			if randVal > g.val {
				mul = g.mul
			}
		}
		if mul == 0 {
			balance = guild.UpdateCoins(uint(amount), "-", userId)
			if balance == nil {
				return fmt.Errorf("Can't find the user.")
			}
			_, err = ctx.Reply(embeds.Info(
				"cOwOins - Gamble",
				fmt.Sprintf("You lost %v cOwOins.\nTry again next time.", amount),
				fmt.Sprintf("You have %v cOwOins currently.", *balance),
			))
		} else {
			balance = guild.UpdateCoins(uint(mul*amount), "+", userId)
			if balance == nil {
				return fmt.Errorf("Can't find the user.")
			}
			_, err = ctx.Reply(embeds.Info(
				"cOwOins - Gamble",
				fmt.Sprintf("UwU **x%v**! You have gained **%v** cOwOins.", mul, mul*amount),
				fmt.Sprintf("You have %v cOwOins currently.", *balance),
			))
		}
		return err
	}
	return fmt.Errorf("Guild not found.")
}
