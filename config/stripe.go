package config

import "github.com/OdyseeTeam/commentron/env"

var connectAPIKey string
var connectAPIKeyTest string

func initStripe(conf *env.Config) {
	connectAPIKey = conf.StripeConnectAPIKey
	connectAPIKeyTest = conf.StripeConnectAPIKeyTest
}

// Environment is a type representing a stripe environment
type Environment string

// StripeTest stripe testing environment
var StripeTest = Environment("test")

// StripeProd stripe production environment
var StripeProd = Environment("live")

// From converts a string to a stripe environment
func From(env string) Environment {
	return Environment(env)
}

// ConnectAPIKey returns the Connect API Key for the passed environment
func ConnectAPIKey(env Environment) string {
	if env == StripeProd {
		return connectAPIKey
	}
	return connectAPIKeyTest
}
