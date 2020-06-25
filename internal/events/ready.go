package events

import (
	"bowot/internal/cron"

	"github.com/andersfylling/disgord"
)

func init() {
	events[disgord.EvtReady] = func(session disgord.Session, evt *disgord.Ready) {
		cron.Start()
	}
}
