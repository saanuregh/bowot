package commands

import (
	"bowot/internal/embeds"
	"fmt"

	"github.com/auttaja/gommand"
	"github.com/bregydoc/gtranslate"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "translate",
		Aliases:     []string{},
		Description: "Translates the given sentence.",
		Category:    utilCategory,
		Function:    translate,
		ArgTransformers: []gommand.ArgTransformer{
			{
				Optional: false,
				Function: gommand.StringTransformer,
			},
			{
				Optional: false,
				Function: gommand.StringTransformer,
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

func translate(ctx *gommand.Context) error {
	from := ctx.Args[0].(string)
	to := ctx.Args[1].(string)
	text := ctx.Args[2].(string)
	translated, err := gtranslate.TranslateWithParams(
		text,
		gtranslate.TranslationParams{
			From:  from,
			To:    to,
			Tries: 3,
		},
	)
	if err != nil {
		return err
	}
	_, err = ctx.Reply(embeds.Info(
		"Translate",
		fmt.Sprintf("**From (%v) : **%v\n**To (%v) : **%v", from, text, to, translated),
		"",
	))
	return err
}
