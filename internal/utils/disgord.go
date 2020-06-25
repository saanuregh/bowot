package utils

import (
	"regexp"

	"github.com/andersfylling/disgord"
)

func GetGuildName(m *disgord.Member) string {
	name := m.Nick
	if name == "" {
		name = m.User.Username
	}
	return name
}

func GetGuildSelfRoles(roles []*disgord.Role, regex string) []string {
	selfroles := []string{}
	for _, r := range roles {
		if regexp.MustCompile(regex).FindString(r.Name) != "" {
			selfroles = append(selfroles, r.ID.String())
		}
	}
	return selfroles
}
