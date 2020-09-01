package env

import (
	"github.com/lbryio/lbry.go/v2/extras/errors"

	e "github.com/caarlos0/env"
)

// Config holds the environment configuration used by lighthouse.
type Config struct {
	MySQLDsn     string `env:"MYSQL_DSN"`
	SlackHookURL string `env:"SLACKHOOKURL"`
	SlackChannel string `env:"SLACKCHANNEL"`
	APIURL       string `env:"APIURL" envDefault:"https://api.lbry.com/event/comment"`
	APIToken     string `env:"APITOKEN"`
}

// NewWithEnvVars creates an Config from environment variables
func NewWithEnvVars() (*Config, error) {
	cfg := &Config{}
	err := e.Parse(cfg)
	if err != nil {
		return nil, errors.Err(err)
	}

	return cfg, nil
}
