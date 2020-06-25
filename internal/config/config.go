package config

import (
	"github.com/kkyr/fig"
)

var C struct {
	Bot struct {
		Token          string `required:"true"`
		DefaultPrefix  string `default:"%"`
		SelfRolePrefix string `default:"color-.*"`
		BotLibre       struct {
			Enabled       bool `default:"true"`
			ApplicationID string
			InstanceID    string
		}
		Owner string `required:"true"`
	}
	Db struct {
		Secret string `required:"true"`
	}
	Logging struct {
		SentryDsn string `required:"true"`
		Level     string `default:"info"`
	}
	Reddit struct {
		Meme      []string `default:"[dankmemes,memes,ComedyCemetery,MemeEconomy,comedyheaven]"`
		CopyPasta []string `default:"[copypasta]"`
		Whoosh    []string `default:"[whoosh]"`
	}
}

func init() {
	err := fig.Load(&C, fig.File("config.yaml"))
	if err != nil {
		panic(err)
	}
}
