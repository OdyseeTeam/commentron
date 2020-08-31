package cmd

import (
	"github.com/lbryio/commentron/meta"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Commentron",
	Long:  `All software has versions. This is Commentron's`,
	Run: func(cmd *cobra.Command, args []string) {
		println("Semantic Version: ", meta.GetSemVersion())
		println("Version: " + meta.GetVersion())
		println("Version(long): " + meta.GetVersionLong())
	},
}
