package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"
)

var (
    MIN_CMD_ARGS = 2 
)

func main() {
    if len(os.Args) < MIN_CMD_ARGS {
	fmt.Println(helpText())
	os.Exit(1)
    }
    
    switch os.Args[1] {
	case "add":
	    proceedAdd()
	case "list":
	    proceedList()
	case "summary":
	    proceedSummary()
	case "delete":
	    proceedDelete()
	default:
	    fmt.Println(helpText())
	    os.Exit(1)
    }
}

type Expence struct {
    Id int	`json:"id"`
    Description string `json:"desc"`
    Amount float64 `json:"amount"`
    CreatedAt int64 `json:"created_at"`
}

func proceedAdd() {
    desc, amount := handleAddFlags()
    
    expencesFile, expences := getExpences()
    defer expencesFile.Close()

    // update content (add new expence) 
    var expenceId int 
    if len(expences) != 0 {
	expenceId = expences[len(expences)-1].Id + 1
    }

    expence := Expence {
	Id: expenceId,
	Description: desc,
	Amount: amount,
	CreatedAt: time.Now().Unix(),
    }
    expences = append(expences, expence)

    // write updated content back
    json, err := json.Marshal(expences)
    if err != nil {
	fmt.Println("Error marshaling file:", err)
	expencesFile.Close()
	os.Exit(1)
    }

    // clean file
    if err = expencesFile.Truncate(0); err != nil {
	fmt.Println("Error truncating file:", err)
	expencesFile.Close()
	os.Exit(1)
    }
    if _, err = expencesFile.Seek(0,0); err != nil {
	fmt.Println("Error seeking file:", err)
	expencesFile.Close()
	os.Exit(1)
    }

    // write updated content to a file
    if _, err = expencesFile.Write(json); err != nil {
	fmt.Println("Error writing to file:", err)
	expencesFile.Close()
	os.Exit(1)
    }

    fmt.Printf("# Expence added successfully (ID: %d)\n", expenceId)   
}

func proceedList() {
    f, expences := getExpences()
    defer f.Close()

    if len(expences) == 0 {
	fmt.Println("No expences yet.")
    } else {
	fmt.Println("# ID\tDate\tDescription\tAmount")
	for _, expence := range expences {
	    created_at := time.Unix(expence.CreatedAt, 0).Format("2006/01/02") 
	    fmt.Printf("# %d\t%s\t%s\t$%.2f\n", expence.Id, created_at, expence.Description, expence.Amount)
	}
    }
}

func proceedSummary() {
    f, expences := getExpences()
    defer f.Close()

    var summary float64
    for _, expence := range expences {
	summary += expence.Amount 
    }
    
    fmt.Printf("# Total expences: $%.2f\n", summary)
}

func proceedDelete() {
    panic("Not implemented yet.")
}

func handleAddFlags() (string, float64) {
    add := flag.NewFlagSet("add", flag.ExitOnError)

    desc := add.String("description", "", "A description for expence.")
    amount := add.Float64("amount", 0, "Expence amount.")

    if err := add.Parse(os.Args[2:]); err != nil {
	fmt.Println("Cannot parse add flags:", err)
	os.Exit(1)
    }

    return *desc, *amount
}

func getExpences() (*os.File, []Expence) {
    // open file
    expencesFile, err := os.OpenFile("expences.json", os.O_RDWR|os.O_CREATE, 0660)   
    if err != nil {
	fmt.Println("Couldn't open expences.json:", err)
	os.Exit(1)
    }

    // get file stats to get length of file
    fileStat, err := expencesFile.Stat()
    if err != nil {
	fmt.Println("Error getting file stat:", err)
	expencesFile.Close()
	os.Exit(1)
    }

    // read content
    var expences []Expence 

    if fileStat.Size() > 0 {
	err = json.NewDecoder(expencesFile).Decode(&expences)
	if err != nil {
	    fmt.Println("Error unmarshaling file:", err)
	    expencesFile.Close()
	    os.Exit(1)
	}
    }

    return expencesFile, expences
}
