package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

func init() {
    rootCmd.AddCommand(commandList)
}

var commandList = &cobra.Command {
    Use: "list",
    Short: "To list all expences.",
    Run: list,
}

func list(cmd *cobra.Command, args []string) {
    f, expences := getExpences()
    defer f.Close()

    if len(expences) == 0 {
	fmt.Println("# No expences yet.")
	if err := os.Remove(expencesFilepath); err != nil {
	    fmt.Println("# Not critical wtf error")
	}
    } else {
	fmt.Println("# ID\tDate\tDescription\tAmount")
	for _, expence := range expences {
	    created_at := time.Unix(expence.CreatedAt, 0).Format("2006/01/02") 
	    fmt.Printf("# %d\t%s\t%s\t$%.2f\n", expence.Id, created_at, expence.Description, expence.Amount)
	}
    }
}
