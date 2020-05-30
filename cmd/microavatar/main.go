package main

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

var logger zerolog.Logger
var verbose bool

var rootCmd = &cobra.Command{
	Use: "microavatar",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout}).Level(zerolog.InfoLevel)
		if verbose {
			logger = logger.Level(zerolog.DebugLevel)
		}
		return nil
	},
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&verbose, "verbose", false, "Verbose logging")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logger.Fatal().Err(err).Msg("Command failed")
	}
}
