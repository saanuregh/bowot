package cron

import (
	"bowot/internal/utils"
	"fmt"
)

const STATUSEXPR = "* * * * *"

func initStatus() {
	statuses := []string{
		"hey",
		"uwu",
		"owo",
	}
	_, err := c.AddFunc(STATUSEXPR, func() {
		err := client.UpdateStatusString(statuses[utils.GetRandomInt(len(statuses))])
		if err != nil {
			client.Logger().Error(fmt.Errorf("CRON STATUS FAILED MSG=%v", err))
			return
		}
		client.Logger().Info("CRON STATUS UPDATED")
	})
	if err != nil {
		client.Logger().Error(fmt.Errorf("CRON STATUS FAILED MSG=%v", err))
	}
	client.Logger().Info("CRON STATUS INITIALIZED")
}
