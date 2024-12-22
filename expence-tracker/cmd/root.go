package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use: "expence-tracker",
    Short: "An expence-tracker is CLI for managing expences.",
}


func Execute() error {
    return rootCmd.Execute()
}
