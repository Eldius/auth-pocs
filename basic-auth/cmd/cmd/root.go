/*
Package cmd is the Basic Auth's subcommands package
*/
package cmd

import (
	"fmt"
	"github.com/eldius/auth-pocs/basic-auth/internal/api"
	"github.com/eldius/auth-pocs/basic-auth/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "basic-auth",
	Short: "A simple Basic Authentication POC",
	Long:  `A simple Basic Authentication POC.`,
	PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
		// TODO init environment configuration
		return config.Setup(cfgFile)
	},
	// Uncomment the following line if your bare application
	// has an action associated with it:
	RunE: func(cmd *cobra.Command, args []string) error {
		// TODO start server
		if err := api.Start(8080); err != nil {
			return err
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var (
	cfgFile string
)

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.basic-auth.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Enable debug log")
	if err := viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug")); err != nil {
		err = fmt.Errorf("binding debug configuration: %w", err)
		panic(err)
	}
}
