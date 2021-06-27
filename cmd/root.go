package cmd

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "wait-for-it",
	Short: "Waiting for service(s)",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// cmd.Help()
		c, _, err := cmd.Find(os.Args[1:])
		// default cmd if no cmd is given
		if err == nil && c.Use == cmd.Use {
			args := append([]string{"execute"}, os.Args[1:]...)
			cmd.SetArgs(args)
		}

		if err := cmd.Execute(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
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
	rootCmd.PersistentFlags().String("level", "info", "Log level")
}

func initConfig() {
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
