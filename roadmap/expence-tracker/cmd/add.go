package cmd

import (
	"fmt"
	"qqweq/playground/roadmap/expence-tracker/models"
	"time"

	"github.com/spf13/cobra"
)

func init() {
    CommandAdd.Flags().StringVar(&desc, "description", "", "A description for expence.")
    CommandAdd.Flags().Float64Var(&amount, "amount", 0.0, "Expence amount.")

    CommandAdd.MarkFlagRequired("description")
    CommandAdd.MarkFlagRequired("amount")

    rootCmd.AddCommand(CommandAdd)
}

var CommandAdd = &cobra.Command{
    Use: "add",
    Short: "To save expence.",
    Run: add,
}

var (
    desc string
    amount float64
)

func add(cmd *cobra.Command, args []string) {
    expences := getExpences()

    // update content (add new expence) 
    var expenceId int 
    if len(expences) != 0 {
	expenceId = expences[len(expences)-1].Id + 1
    }

    expence := models.Expence {
	Id: expenceId,
	Description: desc,
	Amount: amount,
	CreatedAt: time.Now().Unix(),
    }
    expences = append(expences, expence)

    writeExpences(expences)

    fmt.Printf("# Expence added successfully (ID: %d)\n", expenceId)   
}
