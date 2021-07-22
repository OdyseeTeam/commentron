package env

import (
	"github.com/lbryio/lbry.go/v2/extras/errors"

	e "github.com/caarlos0/env"
)

// Config holds the environment configuration used by lighthouse.
type Config struct {
	MySQLDsnRW              string `env:"MYSQL_DSN_RW"`
	MySQLDsnRO              string `env:"MYSQL_DSN_RO"`
	IsTestMode              bool   `env:"IS_TEST"`
	SDKUrl                  string `env:"SDK_URL"`
	SlackHookURL            string `env:"SLACKHOOKURL"`
	SlackChannel            string `env:"SLACKCHANNEL"`
	APIURL                  string `env:"APIURL" envDefault:"https://api.lbry.com/event/comment"`
	APIToken                string `env:"APITOKEN"`
	SocketyToken            string `env:"SOCKETY_TOKEN"`
	TestChannel             string `env:"TEST_CHANNEL"`
	TestURL                 string `env:"TEST_URL" envDefault:"http://localhost:5900/api/v2"`
	StripeConnectAPIKey     string `env:"STRIPE_CONNECT_API_KEY"`
	StripeConnectAPIKeyTest string `env:"STRIPE_CONNECT_API_KEY_TEST"`
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
