package commands

import (
	"bowot/internal/embeds"
	"fmt"
	"time"

	"github.com/auttaja/gommand"
	"github.com/parnurzeal/gorequest"
)

func init() {
	cmds = append(cmds, &gommand.Command{
		Name:        "valorant",
		Aliases:     []string{},
		Description: "Gets valorant asia-pacific server status.",
		Category:    funCategory,
		Function:    valorant,
	})
}

func valorant(ctx *gommand.Context) error {
	type incident struct {
		Description       string      `json:"description"`
		CreatedAt         time.Time   `json:"created_at"`
		Platforms         []string    `json:"platforms"`
		MaintenanceStatus interface{} `json:"maintenance_status"`
		IncidentSeverity  interface{} `json:"incident_severity"`
		Updates           []struct {
			CreatedAt   time.Time `json:"created_at"`
			UpdatedAt   time.Time `json:"updated_at"`
			Description string    `json:"description"`
		} `json:"updates"`
		UpdatedAt interface{} `json:"updated_at"`
	}
	var f []struct {
		Name    string `json:"name"`
		Regions []struct {
			Name         string     `json:"name"`
			Maintenances []incident `json:"maintenances"`
			Incidents    []incident `json:"incidents"`
		} `json:"regions"`
	}
	_, _, errs := gorequest.New().Get("https://riotstatus.vercel.app/valorant").EndStruct(&f)
	if len(errs) > 0 {
		return errs[0]
	}
	data := f[0].Regions[0]
	for _, r := range f[0].Regions {
		if r.Name == "ap" {
			data = r
		}
	}
	incidents := len(data.Incidents) > 0
	maintenances := len(data.Maintenances) > 0
	if !incidents && !maintenances {
		_, err := ctx.Reply(embeds.Info(
			"Valorant Status - Asia",
			"All fine.",
			"",
		))
		return err
	}
	descBuilder := func(incidents []incident) string {
		d := ""
		for i, in := range incidents {
			d += fmt.Sprintf("%v - %v\n", i+1, in.Description)
			d += fmt.Sprintf("Created at: %v\n", in.CreatedAt.Format(time.RFC822))
			if in.UpdatedAt != nil {
				d += fmt.Sprintf("Updated at: %v\n", in.UpdatedAt.(time.Time).Format(time.RFC822))
			}
			if in.IncidentSeverity != nil {
				d += fmt.Sprintf("Severity: %v\n", in.IncidentSeverity)
			}
			if in.MaintenanceStatus != nil {
				d += fmt.Sprintf("Maintenance Status: %v\n", in.MaintenanceStatus)
			}
			if len(in.Updates) > 0 {
				d += "Updates:\n"
				for _, u := range in.Updates {
					d += fmt.Sprintf("> %v - %v\n", u.CreatedAt.Format(time.RFC822), u.Description)
				}
			}
		}
		return d
	}
	desc := ""
	if incidents {
		desc += fmt.Sprintf("__**Incidents - %v**__\n%v", len(data.Incidents), descBuilder(data.Incidents))
	}
	if maintenances {
		desc += fmt.Sprintf("__**Maintenance - %v**__\n%v", len(data.Maintenances), descBuilder(data.Maintenances))
	}
	_, err := ctx.Reply(embeds.Info(
		"Valorant Status - Asia",
		desc,
		"",
	))
	return err
}
