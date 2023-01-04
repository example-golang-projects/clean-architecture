package main

import (
	"context"
	masterCfg "github.com/example-golang-projects/clean-architecture/cmd/server/master/config"
	"github.com/example-golang-projects/clean-architecture/cmd/server/user"
	userCfg "github.com/example-golang-projects/clean-architecture/cmd/server/user/config"
	"github.com/example-golang-projects/clean-architecture/golibs/config"
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{}

func main() {
	makeRootCmd()
	if err := rootCmd.ExecuteContext(context.Background()); err != nil {
		os.Exit(1)
	}
}

func makeRootCmd() {
	var ()

	userCmd := &cobra.Command{
		Use:   "user",
		Short: "Start User server",
		Long:  `User is internal back-end to manage user, role, permission domains)`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			userConfig := userCfg.Config{}
			err := config.MustLoad(config.FileType_JSON, "./cmd/development/secret/config/user/config.local.json", &userConfig)
			if err != nil {
				panic(err)
			}
			user.RunUserService(userConfig)
		},
	}
	masterCmd := &cobra.Command{
		Use:   "master",
		Short: "Start Master server",
		Long:  `User is internal back-end to manage master data (import, export data which hardly change)`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			masterCfg := &masterCfg.Config{}
			err := config.MustLoad(config.FileType_JSON, "./cmd/development/secret/config/master/config.local.json", masterCfg)
			if err != nil {
				panic(err)
			}
		},
	}

	rootCmd.AddCommand(
		userCmd,
		masterCmd,
	)
}
