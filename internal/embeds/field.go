package embeds

import "github.com/andersfylling/disgord"

// Field instantiates an embed field.
func Field(name, value string, inline bool) *disgord.EmbedField {
	return &disgord.EmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	}
}
