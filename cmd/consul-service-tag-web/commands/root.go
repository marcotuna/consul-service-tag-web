package commands

import (
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Command ...
type Command = cobra.Command

// Run ...
func Run(args []string) error {
	RootCmd.SetArgs(args)
	return RootCmd.Execute()
}

// RootCmd ...
var RootCmd = &cobra.Command{
	Use:   "consul-service-tag-web",
	Short: "Webui for Consul Service Tags",
	Long:  `For DevOPS environments where multiple teams need to get information related to current running services`,
}

func init() {
	RootCmd.PersistentFlags().StringP("tag.filter", "t", "", "Filter tag to use.")

	viper.SetEnvPrefix("cstw")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	viper.BindPFlag("tag-filter", RootCmd.PersistentFlags().Lookup("tag.filter"))
}
