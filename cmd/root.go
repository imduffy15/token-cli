package cmd

import (
	"fmt"
	"os"

	"github.com/imduffy15/token-cli/cli"
	"github.com/imduffy15/token-cli/client"
	"github.com/imduffy15/token-cli/config"
	"github.com/imduffy15/token-cli/help"
	"github.com/imduffy15/token-cli/version"
	"github.com/spf13/cobra"
)

var cfgFile client.Config
var log cli.Logger

// Global flags
var (
	skipSSLValidation bool
	verbose           bool
)

var (
	scope string
	port  int
	force bool
)

var RootCmd = cobra.Command{
	Use:   "token-cli",
	Short: "A cli for generating tokens",
	Long:  help.Root(version.VersionString()),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "See additional info on HTTP requests")
	log = cli.NewLogger(os.Stderr, os.Stdout, os.Stderr, os.Stderr)
}

func initConfig() {
	// Startup tasks
}

func GetLogger() *cli.Logger {
	return &log
}

func GetSavedConfig() client.Config {
	cfgFile = config.Read()
	cfgFile.Verbose = verbose
	return cfgFile
}
