package env

import (
	e "github.com/caarlos0/env"
	"github.com/lbryio/lbry.go/v2/extras/errors"
)

// Config holds the environment configuration used by lighthouse.
type Config struct {
	MySQLDsnRW                        string `env:"MYSQL_DSN_RW"`
	MySQLDsnRO                        string `env:"MYSQL_DSN_RO"`
	IsTestMode                        bool   `env:"IS_TEST"`
	SDKUrl                            string `env:"SDK_URL"`
	SlackHookURL                      string `env:"SLACKHOOKURL"`
	SlackChannel                      string `env:"SLACKCHANNEL"`
	APIURL                            string `env:"APIURL" envDefault:"https://api.odysee.com"`
	APIToken                          string `env:"APITOKEN"`
	SocketyToken                      string `env:"SOCKETY_TOKEN"`
	TestChannel                       string `env:"TEST_CHANNEL"`
	TestURL                           string `env:"TEST_URL" envDefault:"http://localhost:5900/api/v2"`
	StripeConnectAPIKey               string `env:"STRIPE_CONNECT_API_KEY"`
	StripeConnectAPIKeyTest           string `env:"STRIPE_CONNECT_API_KEY_TEST"`
	CommentClassifierAPIURL           string `env:"COMMENT_CLASSIFIER_API_URL" envDefault:"https://localhost/_/comments"`
	CommentClassificationEndpointAuth string `env:"COMMENT_CLASSIFICATION_ENDPOINT_AUTH" envDefault:""`
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
