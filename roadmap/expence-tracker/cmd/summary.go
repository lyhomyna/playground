package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

func init() {
    commandSummary.Flags().Int("month", -1, "Display amount per month. (Numeric values represent months: 1,2,..12)")

    rootCmd.AddCommand(commandSummary)
}

var commandSummary = &cobra.Command { 
    Use: "summary",
    Short: "To summarize expences amount.",
    Run: summary,
}

func summary(cmd *cobra.Command, args []string) {
    month, _ := cmd.Flags().GetInt("month")
    if month != -1 && (month < 1 || month > 12) {
	_ = cmd.Help()
	os.Exit(1)
    }

    expences := getExpences()
    if month == -1 {
	var summary float64
	for _, expence := range expences {
	    summary += expence.Amount 
	}
	
	fmt.Printf("# Total expences: $%.2f\n", summary)
    } else {
	var monthSummary float64
	for _, expence := range expences {
	    date := time.Unix(expence.CreatedAt, 0).Format("1")

	    expenceMonth, _ := strconv.ParseInt(date, 10, 0)

	    if int(expenceMonth) == month {
		monthSummary += expence.Amount 
	    }
	}

	fmt.Printf("# Expences by %d month: $%.2f\n", month, monthSummary)
    }
}
