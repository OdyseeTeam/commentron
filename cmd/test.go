package cmd

import (
	"github.com/OdyseeTeam/commentron/config"
	"github.com/OdyseeTeam/commentron/env"
	"github.com/OdyseeTeam/commentron/server/lbry"
	"github.com/OdyseeTeam/commentron/tests"

	"github.com/pkg/profile"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Runs the tests of commentron against an instance",
	Long:  `Runs the tests of commentron against an instance`,
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("codeprofile") {
			defer profile.Start(profile.NoShutdownHook).Stop()
		}
		conf, err := env.NewWithEnvVars()
		if err != nil {
			logrus.Panic(err)
		}
		config.InitializeConfiguration(conf)
		lbry.Init(conf)
		tests.Launch()
	},
}
