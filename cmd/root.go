// Copyright © 2021 Bitrock s.r.l. <devops@bitrock.it>
package cmd

import (
	"fmt"
	"io"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile  string
	logLevel string
	jsonLogs bool
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "caravan",
	Short: "This tool allow you to deploy and setup a caravan cluster",
	Long: `With caravan cli you can:
		- init the required provider
		- bake the corresponding images
		- deploy the infrastructure on the given provider with terraform
		- destroy the infrastructure
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
	SilenceUsage:      true,
	PersistentPreRunE: rootPreRun,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.caravan.yaml)")
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "log level to be used")
	rootCmd.PersistentFlags().BoolVar(&jsonLogs, "json-logs", false, "log in JSON format")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".caravan" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".caravan")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}

func setUpLogs(out io.Writer, level zerolog.Level, humanLogs bool) error {
	zerolog.SetGlobalLevel(level)
	if humanLogs {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: out})
	}
	return nil
}

func rootPreRun(cmd *cobra.Command, args []string) error {
	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		return fmt.Errorf("%s is not a valid log level: %w", logLevel, err)
	}

	if err = setUpLogs(os.Stdout, level, !jsonLogs); err != nil {
		return err
	}

	log.Info().Msg("Logger initialized")
	return nil
}
