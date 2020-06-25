package commands

import (
	"bowot/internal/embeds"
	"fmt"
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/auttaja/gommand"
	"github.com/parnurzeal/gorequest"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:            "define",
		Aliases:         []string{"meaning"},
		Description:     "Defines the given word.",
		Category:        utilCategory,
		Function:        define,
		ArgTransformers: []gommand.ArgTransformer{{Optional: false, Function: gommand.StringTransformer}},
	})
}

func define(ctx *gommand.Context) error {
	word := ctx.Args[0].(string)
	var f []struct {
		Word     string                              `json:"word"`
		Phonetic string                              `json:"phonetic"`
		Origin   string                              `json:"origin,omitempty"`
		Meaning  map[string][]map[string]interface{} `json:"meaning"`
	}
	_, _, errs := gorequest.New().Get("https://googledictionaryapi.eu-gb.mybluemix.net").Param("define", word).EndStruct(&f)
	if len(errs) > 0 {
		ctx.Session.Logger().Error(errs[0])
		return fmt.Errorf("Something happned, a glitch in the matrix!")
	}
	embs := make([]*disgord.EmbedField, 0)
	embs = append(embs, &disgord.EmbedField{Name: "Word", Value: f[0].Word, Inline: true})
	embs = append(embs, &disgord.EmbedField{Name: "Phonetic", Value: f[0].Phonetic, Inline: true})
	if f[0].Origin != "" {
		embs = append(embs, &disgord.EmbedField{Name: "Origin", Value: f[0].Origin, Inline: false})
	}
	s := ""
	i := 0
	for k1, v1 := range f[0].Meaning {
		i += 1
		s = s + fmt.Sprintf("%d. **%s**\n", i, k1)
		for _, v2 := range v1 {
			s = s + fmt.Sprintf("**Definition**: %s\n", v2["definition"].(string))
			s = s + fmt.Sprintf("**Example**: %s\n", v2["example"].(string))
			if _syns, ok := v2["synonyms"].([]interface{}); ok {
				syns := make([]string, len(_syns))
				for i, v := range _syns {
					syns[i] = fmt.Sprint(v)
				}
				s = s + fmt.Sprintf("**Synonyms**: %s\n", strings.Join(syns, ", "))
			}
			s = s + "\n"
		}
	}
	embs = append(embs, &disgord.EmbedField{Name: "Meaning", Value: s, Inline: false})
	_, err := ctx.Reply(embeds.Info(
		"Define",
		"",
		"",
		embs...,
	))
	return err
}
