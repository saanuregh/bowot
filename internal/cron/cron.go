package cron

import (
	"github.com/andersfylling/disgord"
	"github.com/robfig/cron/v3"
)

var client *disgord.Client
var c *cron.Cron

func Register(disgordClient *disgord.Client) {
	client = disgordClient
	c = cron.New()
	initHydrate()
	initStatus()
}

func Start() {
	c.Start()
	client.Logger().Info("CRON STARTED")
}

func Stop() {
	c.Stop()
	client.Logger().Info("CRON STOPED")
}
