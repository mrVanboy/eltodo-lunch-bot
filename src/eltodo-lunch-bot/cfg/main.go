package cfg

import (
	"github.com/caarlos0/env"
	"errors"
)

type config struct {
	PathToPdfToText	string	`env:"PATH_PDFTOTEXT" envDefault:"/usr/bin/pdftotext"`
	TimeZone		string 	`env:"TIMEZONE" envDefault:"Europe/Prague"`
	Cron			string 	`env:"CRON" envDefault:"0 10 * * 1-5"`
	Footer			string 	`env:"FOOTER" envDefault:"Neručíme za věrohodnost údajů. Vždy zkontrolujte oficiální nabídku."`
	BotName			string 	`env:"BOT_NAME" envDefault:"obědbot"`
	MainWebhookUrl	string	`env:"WEBHOOK_MAIN_URL"`
	DebugWebhookUrl	string	`env:"WEBHOOK_DEBUG_URL"`

	UrlBK			string 	`env:"URL_BK"`
	UrlDC			string 	`env:"URL_DC"`
	UrlNK			string 	`env:"URL_NK"`
	UrlPP			string 	`env:"URL_PP"`
}

var c *config

func Load() error {
	c = new(config)
	err := env.Parse(c)
	if err != nil {
		return err
	}
	if c.MainWebhookUrl == `` {
		return errors.New(`env WEBHOOK_MAIN_URL must be specified`)
	}
	if c.DebugWebhookUrl == `` {
		return errors.New(`env WEBHOOK_DEBUG_URL must be specified`)
	}
	return nil
}

func Get() *config {
	return c
}