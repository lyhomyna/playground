package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(commandSummary)
}

var commandSummary = &cobra.Command { 
    Use: "summary",
    Short: "To summarize expences amount.",
    Run: summary,
}

func summary(cmd *cobra.Command, args []string) {
    f, expences := getExpences()
    defer f.Close()

    var summary float64
    for _, expence := range expences {
	summary += expence.Amount 
    }
    
    fmt.Printf("# Total expences: $%.2f\n", summary)
}
