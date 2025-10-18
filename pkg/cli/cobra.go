package cli

import (
    "github.com/spf13/cobra"
    "os"
)

var rootCmd = &cobra.Command{
    Use:   "artemis",
    Short: "Artemis CLI - Build monoliths, modularly",
    Long:  `Artemis is a modular framework for Go that helps you build monolithic applications with a structured approach.`,
    Run: func(cmd *cobra.Command, args []string) {
        cmd.Help()
    },
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        os.Exit(1)
    }
}

func init() {
    // Here you can add subcommands to the root command
    // For example:
    // rootCmd.AddCommand(newCmd)
    // rootCmd.AddCommand(makeCmd)
    // rootCmd.AddCommand(versionCmd)
}