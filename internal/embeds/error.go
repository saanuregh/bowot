package embeds

import (
	"github.com/andersfylling/disgord"
)

// Error instantiates an embed to be used to throw errors.
func Error(title string, err error, report bool) *disgord.Embed {
	footer := "This is a user error."
	if report {
		footer = "This has been reported to the developers."
	}
	return &disgord.Embed{
		Title:       title,
		Description: err.Error(),
		Footer:      &disgord.EmbedFooter{Text: footer},
		Color:       0xf76580,
	}
}

// ErrorString instantiates an embed to be used to throw errors as strings.
func ErrorString(title, message string) *disgord.Embed {
	return &disgord.Embed{
		Title:       title,
		Description: message,
		Color:       0xf76580,
	}
}
