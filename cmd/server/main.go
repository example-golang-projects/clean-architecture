package main

import (
	"context"
	"fmt"
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
			fmt.Println("test user na")
			// Set config here
		},
	}
	masterCmd := &cobra.Command{
		Use:   "master",
		Short: "Start Master server",
		Long:  `User is internal back-end to manage master data`,
		Args:  cobra.MinimumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("test master na")
			// Set config here
		},
	}

	rootCmd.AddCommand(
		userCmd,
		masterCmd,
	)
}
