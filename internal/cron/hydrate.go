package cron

import (
	"bowot/internal/db"
	"context"
	"fmt"

	"github.com/andersfylling/disgord"
	"github.com/lrita/cmap"
	"github.com/robfig/cron/v3"
)

var HydrateUsers cmap.Cmap

const HYDRATEEXPR = "0 * * * *"

func (h *HydrateUser) sentHydrateDM() {
	u, err := client.GetUser(context.Background(), disgord.ParseSnowflakeString(h.userID))
	if err != nil {
		client.Logger().Error(fmt.Errorf("CRON HYDRATE DM FAILED USERID=%s MSG=%v", h.userID, err))
		return
	}
	if h.status == "online" {
		_, _, err = u.SendMsgString(context.Background(), client, "drink water pls")
		if err != nil {
			client.Logger().Error(fmt.Errorf("CRON HYDRATE DM FAILED USERID=%s MSG=%v", h.userID, err))
			return
		}
		client.Logger().Info(fmt.Sprintf("CRON HYDRATE DM SUCCESS USERID=%s", h.userID))
	}
}

type HydrateUser struct {
	userID string
	status string
	cronID cron.EntryID
}

func initHydrate() {
	client.On(disgord.EvtPresenceUpdate, func(session disgord.Session, evt *disgord.PresenceUpdate) {
		HydrateUsers.Range(func(key, value interface{}) bool {
			if key.(string) == evt.User.ID.String() {
				value.(*HydrateUser).status = evt.Status
			}
			return true
		})
	})
	db.Guilds.Range(
		func(key, value interface{}) bool {
			g := value.(*db.Guild)
			guild, err := client.GetGuild(context.Background(), disgord.ParseSnowflakeString(g.ID))
			if err != nil {
				client.Logger().Error(err)
				return true
			}
			for _, h := range g.Hydrate {
				if _, ok := HydrateUsers.Load(h); !ok {
					status := "online"
					for _, p := range guild.Presences {
						if p.User.ID == disgord.ParseSnowflakeString(h) {
							status = p.Status
						}
					}
					id, err := c.AddFunc(HYDRATEEXPR, func() {
						if tmp, ok := HydrateUsers.Load(h); ok {
							tmp.(*HydrateUser).sentHydrateDM()
						}
					})
					if err != nil {
						client.Logger().Error(err)
						return true
					}
					HydrateUsers.Store(h, &HydrateUser{status: status, cronID: id, userID: h})
				}
			}
			return true
		},
	)
	client.Logger().Info("CRON HYDRATE INITIALIZED")
}

func AddHydrate(userID, guildID string) error {
	if tmp, ok := db.Guilds.Load(guildID); ok {
		g := tmp.(*db.Guild)
		h := g.AddHydrate(userID)
		if h == nil {
			return fmt.Errorf("Couldn't add the user.")
		}
		guild, err := client.GetGuild(context.Background(), disgord.ParseSnowflakeString(guildID))
		if err != nil {
			return err
		}
		status := "offline"
		for _, p := range guild.Presences {
			if p.User.ID == disgord.ParseSnowflakeString(userID) {
				status = p.Status
			}
		}
		id, err := c.AddFunc(HYDRATEEXPR, func() {
			if tmp, ok := HydrateUsers.Load(userID); ok {
				tmp.(*HydrateUser).sentHydrateDM()
			}
		})
		if err != nil {
			return err
		}
		HydrateUsers.Store(h, &HydrateUser{status: status, cronID: id, userID: userID})
		return nil
	}
	return fmt.Errorf("Couldn't find the guild.")
}

func RemoveHydrate(userID, guildID string) error {
	if tmp, ok := db.Guilds.Load(guildID); ok {
		g := tmp.(*db.Guild)
		h := g.RemoveHydrate(userID)
		if h == nil {
			return fmt.Errorf("cant remove hydrate")
		}
		if tmp, ok := HydrateUsers.Load(userID); ok {
			huser := tmp.(*HydrateUser)
			c.Remove(huser.cronID)
			HydrateUsers.Delete(userID)
		}
		return nil
	}
	return fmt.Errorf("Couldn't find the guild.")
}
