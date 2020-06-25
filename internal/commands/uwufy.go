package commands

import (
	"strings"

	"github.com/auttaja/gommand"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "uwufy",
		Aliases:     []string{},
		Description: "uwufy.",
		Category:    funCategory,
		Function:    uwufy,
		ArgTransformers: []gommand.ArgTransformer{
			{
				Optional:  false,
				Function:  gommand.StringTransformer,
				Greedy:    true,
				Remainder: true,
			},
		},
	})
}

func uwufy(ctx *gommand.Context) error {
	str := ctx.Args[0].(string)
	str = strings.ReplaceAll(str, "r", "w")
	str = strings.ReplaceAll(str, "l", "w")
	_, err := ctx.Reply(str)
	return err
}
