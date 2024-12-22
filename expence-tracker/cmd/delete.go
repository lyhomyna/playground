package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
    CommandDel.Flags().IntVar(&id, "id", -1, "Expence id")

    rootCmd.AddCommand(CommandDel)
}

var CommandDel = &cobra.Command{
    Use: "delete",
    Short: "To delete expences.",
    Run: del,
}

var id int 

func del(cmd *cobra.Command, args []string) {
    expences := getExpences()

    if len(expences) == 0 || id == -1 {
	// delete all
	os.Remove(expencesFilepath)
    } else {
	for i, expence := range expences {
	    if expence.Id == id {
		expences = append(expences[:i], expences[i+1:]...)
		break
	    }
	}
	writeExpences(expences)
    }

    fmt.Println("# Expences deleted successfully")   
}
