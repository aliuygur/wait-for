package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wait-for",
	Short: "Test environments",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Verbose (log.level = trace)")
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "Debug (log.level = debug)")
	rootCmd.PersistentFlags().BoolP("execute", "e", false, "Execute")
	rootCmd.PersistentFlags().String("level", "info", "Log level")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".wait-for" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".wait-for")
	}

	viper.SetEnvPrefix("ENVT")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	initLogging()
}

func initLogging() {
	levelStr, _ := rootCmd.PersistentFlags().GetString("level")
	level, err := log.ParseLevel(levelStr)
	if err != nil {
		level = log.InfoLevel
	}

	verbose, _ := rootCmd.PersistentFlags().GetBool("verbose")
	debug, _ := rootCmd.PersistentFlags().GetBool("debug")

	if verbose {
		level = log.TraceLevel
	} else if debug {
		level = log.DebugLevel
	}

	if verbose || debug {
		log.Printf("log.SetLevel(%s)", level)
	}

	log.SetLevel(level)
}
