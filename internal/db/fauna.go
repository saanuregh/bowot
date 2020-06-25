package db

import (
	"bowot/internal/config"
	"bowot/internal/logger"
	"bowot/internal/utils"
	"fmt"
	"regexp"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/lrita/cmap"
	f "github.com/mbraeutig/faunadb-v2/faunadb"
)

var Guilds cmap.Cmap

type Guild struct {
	ref             f.Value    `fauna:"-"`
	ID              string     `fauna:"id"`
	SelfRoles       []string   `fauna:"selfroles"`
	Members         []Member   `fauna:"members"`
	Prefix          string     `fauna:"prefix"`
	SelfRolesRegex  string     `fauna:"selfrolesregex"`
	DefaultSelfRole string     `fauna:"selfrolesregex"`
	CustomCommands  [][]string `fauna:"customcommands"`
	WakePhrases     [][]string `fauna:"wakephrases"`
	Hydrate         []string   `fauna:"hydrate"`
}

type Member struct {
	ID        string `fauna:"id"`
	SelfRole  string `fauna:"selfrole"`
	Coins     uint   `fauna:"coins"`
	LastDaily int64  `fauna:"last_daily"`
}

var client *f.FaunaClient

func init() {
	client = f.NewFaunaClient(config.C.Db.Secret)
	result, err := client.Query(f.Paginate(f.Documents(f.Collection("guilds"))))
	if err != nil {
		logger.CustomLogger.Panic(err)
	}
	var data f.ArrayV
	err = result.At(f.ObjKey("data")).Get(&data)
	if err != nil {
		logger.CustomLogger.Panic(err)
	}
	for _, v := range data {
		var va f.RefV
		v.Get(&va)
		result, err := client.Query(f.Get(va))
		if err != nil {
			logger.CustomLogger.Panic(err)
		}
		guild := Guild{}
		err = result.At(f.ObjKey("data")).Get(&guild)
		if err != nil {
			logger.CustomLogger.Panic(err)
		}
		guild.ref = va
		Guilds.Store(guild.ID, &guild)
	}
	logger.CustomLogger.Info("DB INITIALIZING")
}

func AddGuild(g *disgord.Guild, ms []*disgord.Member) *Guild {
	selfroles := utils.GetGuildSelfRoles(g.Roles, config.C.Bot.SelfRolePrefix)
	guild := Guild{
		ID:              g.ID.String(),
		Members:         make([]Member, 0),
		SelfRoles:       selfroles,
		Prefix:          config.C.Bot.DefaultPrefix,
		SelfRolesRegex:  config.C.Bot.SelfRolePrefix,
		CustomCommands:  make([][]string, 0),
		WakePhrases:     [][]string{{"owo", "OwO"}, {"uwu", "UwU"}, {"hi", "hi", "👋"}},
		Hydrate:         make([]string, 0),
		DefaultSelfRole: "",
	}
	for _, m := range ms {
		if m.User.Bot {
			continue
		}
		guild.Members = append(guild.Members, Member{
			ID:        m.User.ID.String(),
			Coins:     0,
			LastDaily: time.Now().AddDate(0, 0, -1).UnixNano(),
			SelfRole:  "",
		})
	}
	result, err := client.Query(f.Create(f.Collection("guilds"), f.Obj{"data": guild}))
	if err != nil {
		logger.CustomLogger.Error(err)
		return nil
	}
	ref, err := result.At(f.ObjKey("ref")).GetValue()
	if err != nil {
		logger.CustomLogger.Error(err)
		return nil
	}
	guild.ref = ref
	Guilds.Store(guild.ID, &guild)
	logger.CustomLogger.Info(fmt.Sprintf("DB ADD GUILD ID=%s", g.ID))
	return &guild
}

func RemoveGuild(guildID string) *Guild {
	if tmp, ok := Guilds.Load(guildID); ok {
		guild := tmp.(*Guild)
		_, err := client.Query(f.Delete(guild.ref))
		if err != nil {
			return nil
		}
		Guilds.Delete(guildID)
		logger.CustomLogger.Info(fmt.Sprintf("DB REMOVE GUILD ID=%s", guildID))
		return guild
	}
	return nil
}

func (g *Guild) Sync() error {
	_, err := client.Query(f.Update(g.ref, f.Obj{"data": g}))
	if err != nil {
		logger.CustomLogger.Error(fmt.Errorf("DB SYNCING FAILED MSG=%v", err))
		return err
	}
	logger.CustomLogger.Info("DB SYNCING SUCCESS")
	return nil
}

func (g *Guild) GetMember(memberID string) (*Member, int) {
	for i, member := range g.Members {
		if memberID == member.ID {
			return &member, i
		}
	}
	return nil, -1
}

func (g *Guild) AddMember(m *disgord.Member) *Member {
	member := Member{
		ID:        m.User.ID.String(),
		Coins:     0,
		LastDaily: time.Now().AddDate(0, 0, -1).UnixNano(),
		SelfRole:  g.DefaultSelfRole,
	}
	g.Members = append(g.Members, member)
	if err := g.Sync(); err != nil {
		return nil
	}
	logger.CustomLogger.Info(fmt.Sprintf("DB ADD MEMBER GUILDID=%s MEMBERID=%s", g.ID, m.User.ID))
	return &member
}

func (g *Guild) RemoveMember(memberID string) *Member {
	member, i := g.GetMember(memberID)
	if i == -1 {
		return nil
	}
	g.Members = append(g.Members[:i], g.Members[i+1:]...)
	if err := g.Sync(); err != nil {
		return nil
	}
	logger.CustomLogger.Info(fmt.Sprintf("DB REMOVE MEMBER GUILDID=%s MEMBERID=%s", g.ID, memberID))
	return member
}

func (g *Guild) GetSelfRole(memberID string) *string {
	member, _ := g.GetMember(memberID)
	if member == nil {
		return nil
	}
	return &member.SelfRole
}

func (g *Guild) UpdateSelfRole(selfrole string, memberID string) *string {
	_, i := g.GetMember(memberID)
	if i == -1 {
		return nil
	}
	g.Members[i].SelfRole = selfrole
	if err := g.Sync(); err != nil {
		return nil
	}
	logger.CustomLogger.Info(fmt.Sprintf("DB UPDATE SELFROLE GUILDID=%s MEMBERID=%s ROLEID=%s", g.ID, memberID, selfrole))
	return &g.Members[i].SelfRole
}

func (g *Guild) GetCoins(memberID string) *uint {
	member, _ := g.GetMember(memberID)
	if member == nil {
		return nil
	}
	return &member.Coins
}

func (g *Guild) UpdateCoins(amount uint, op string, memberID string) *uint {
	_, i := g.GetMember(memberID)
	if i == -1 {
		return nil
	}
	if op == "+" {
		g.Members[i].Coins += amount
	}
	if op == "-" {
		g.Members[i].Coins -= amount
	}
	if err := g.Sync(); err != nil {
		return nil
	}
	logger.CustomLogger.Info(fmt.Sprintf("DB UPDATE COIN GUILDID=%s MEMBERID=%s OP=%s AMOUNT=%d", g.ID, memberID, op, amount))
	return &g.Members[i].Coins
}

func (g *Guild) GetLastDaily(memberID string) *int64 {
	_, i := g.GetMember(memberID)
	if i == -1 {
		return nil
	}
	return &g.Members[i].LastDaily
}

func (g *Guild) UpdateLastDaily(memberID string) *int64 {
	_, i := g.GetMember(memberID)
	if i == -1 {
		return nil
	}
	g.Members[i].LastDaily = time.Now().UnixNano()
	if err := g.Sync(); err != nil {
		return nil
	}
	logger.CustomLogger.Info(fmt.Sprintf("DB UPDATE LASTDAILY GUILDID=%s MEMBERID=%s", g.ID, memberID))
	return &g.Members[i].LastDaily
}

func (g *Guild) UpdatePrefix(prefix string) *string {
	g.Prefix = prefix
	if err := g.Sync(); err != nil {
		return nil
	}
	logger.CustomLogger.Info(fmt.Sprintf("DB UPDATE PREFIX ID=%s PREFIX=%s", g.ID, prefix))
	return &g.Prefix
}

func (g *Guild) UpdateDefaultSelfRole(id string) *string {
	g.DefaultSelfRole = id
	if err := g.Sync(); err != nil {
		return nil
	}
	logger.CustomLogger.Info(fmt.Sprintf("DB UPDATE DEFAULT SELFROLE ID=%s SELFROLE=%s", g.ID, id))
	return &g.DefaultSelfRole
}

func (g *Guild) UpdateSelfRolesRegex(regex string) *string {
	_, err := regexp.Compile(regex)
	if err != nil {
		return nil
	}
	g.SelfRolesRegex = regex
	if err := g.Sync(); err != nil {
		return nil
	}
	logger.CustomLogger.Info(fmt.Sprintf("DB UPDATE SELFROLE REGEX ID=%s SELFROLE=%s", g.ID, regex))
	return &g.SelfRolesRegex
}

func (g *Guild) UpdateSelfRoles(selfroles []string) *[]string {
	g.SelfRoles = selfroles
	flag := true
	for _, s := range selfroles {
		if s == g.DefaultSelfRole {
			flag = false
		}
	}
	if flag {
		g.DefaultSelfRole = ""
	}
	if len(selfroles) > 0 {
		for i, m := range g.Members {
			flag := true
			for _, r := range selfroles {
				if m.SelfRole == r {
					flag = false
				}
			}
			if flag {
				g.Members[i].SelfRole = g.DefaultSelfRole
			}
		}
	}
	if err := g.Sync(); err != nil {
		return nil
	}
	logger.CustomLogger.Info(fmt.Sprintf("DB UPDATE SELFROLES ID=%s N_SELFROLE=%d", g.ID, len(selfroles)))
	return &g.SelfRoles
}

func (g *Guild) AddCustomCommands(cmd []string) *[][]string {
	g.CustomCommands = append(g.CustomCommands, cmd)
	if err := g.Sync(); err != nil {
		return nil
	}
	logger.CustomLogger.Info(fmt.Sprintf("DB ADD CUSTOMCOMMANDS ID=%s COMMAND=%s REPLY=%s", g.ID, cmd[0], cmd[1]))
	return &g.CustomCommands
}

func (g *Guild) RemoveCustomCommands(cmd string) *[][]string {
	for i, c := range g.CustomCommands {
		if c[0] == cmd {
			g.CustomCommands = append(g.CustomCommands[:i], g.CustomCommands[i+1:]...)
			if err := g.Sync(); err != nil {
				return nil
			}
			logger.CustomLogger.Info(fmt.Sprintf("DB REMOVE CUSTOMCOMMANDS ID=%s COMMAND=%s", g.ID, cmd))
			return &g.CustomCommands
		}
	}
	return nil
}

func (g *Guild) AddWakePhrase(wp []string) *[][]string {
	g.WakePhrases = append(g.WakePhrases, wp)
	if err := g.Sync(); err != nil {
		return nil
	}
	logger.CustomLogger.Info(fmt.Sprintf("DB ADD WAKEPHRASE ID=%s PHRASE=%s REPLY=%s", g.ID, wp[0], wp[1]))
	return &g.WakePhrases
}

func (g *Guild) RemoveWakePhrase(wp string) *[][]string {
	for i, c := range g.WakePhrases {
		if c[0] == wp {
			g.WakePhrases = append(g.WakePhrases[:i], g.WakePhrases[i+1:]...)
			if err := g.Sync(); err != nil {
				return nil
			}
			logger.CustomLogger.Info(fmt.Sprintf("DB REMOVE WAKEPHRASE ID=%s PHRASE=%s", g.ID, wp))
			return &g.WakePhrases
		}
	}
	return nil
}

func (g *Guild) AddHydrate(userID string) *[]string {
	g.Hydrate = append(g.Hydrate, userID)
	if err := g.Sync(); err != nil {
		return nil
	}
	logger.CustomLogger.Info(fmt.Sprintf("DB ADD HYDRATE ID=%s USERID=%s", g.ID, userID))
	return &g.Hydrate
}

func (g *Guild) RemoveHydrate(userID string) *[]string {
	for i, c := range g.Hydrate {
		if c == userID {
			g.Hydrate = append(g.Hydrate[:i], g.Hydrate[i+1:]...)
			if err := g.Sync(); err != nil {
				return nil
			}
			logger.CustomLogger.Info(fmt.Sprintf("DB REMOVE HYDRATE ID=%s USERID=%s", g.ID, userID))
			return &g.Hydrate
		}
	}
	return nil
}
