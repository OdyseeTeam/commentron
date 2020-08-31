package cmd

import (
	"github.com/lbryio/commentron/config"
	"github.com/lbryio/commentron/server"
	"github.com/pkg/profile"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	serveCmd.PersistentFlags().StringVarP(&server.RPCHost, "host", "", "", "host to listen on")
	serveCmd.PersistentFlags().IntVarP(&server.RPCPort, "port", "p", 5900, "port binding used for the rpc server")
	//Bind to Viper
	err := viper.BindPFlags(serveCmd.PersistentFlags())
	if err != nil {
		logrus.Panic(err)
	}
	rootCmd.AddCommand(serveCmd)
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Runs the Commentron JSON RPC server",
	Long:  `Runs the Commentron JSON RPC server`,
	Args:  cobra.OnlyValidArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if viper.GetBool("codeprofile") {
			defer profile.Start(profile.NoShutdownHook).Stop()
		}
		config.InitializeConfiguration()
		server.Start()
	},
}
