package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
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
    Amount int `json:"amount"`
}

func proceedAdd() {
    desc, amount := handleAddFlags()

    // open file
    expencesFile, err := os.OpenFile("expences.json", os.O_RDWR|os.O_CREATE, 0660)   
    if err != nil {
	fmt.Println("Couldn't open expences.json:", err)
	os.Exit(1)
    }
    defer expencesFile.Close()

    // get file stat to get length of file
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

    // add new expence 
    var expenceId int 
    if len(expences) != 0 {
	expenceId = expences[len(expences)-1].Id + 1
    }

    expence := Expence {
	Id: expenceId,
	Description: desc,
	Amount: amount,
    }

    expences = append(expences, expence)

    fmt.Println(expences)

    // write updated content back
    err = json.NewEncoder(expencesFile).Encode(expences)
    if err != nil {
	fmt.Println("Encoding error:", err)
	expencesFile.Close()
	os.Exit(1)
    }

    fmt.Printf("# Expence added successfully (ID: %d)\n", expenceId)   
}

func proceedList() {
    panic("Not implemented yet.")
}

func proceedSummary() {
    panic("Not implemented yet.")
}

func proceedDelete() {
    panic("Not implemented yet.")
}

func handleAddFlags() (string, int) {
    add := flag.NewFlagSet("add", flag.ExitOnError)

    desc := add.String("description", "", "A description for expence.")
    amount := add.Int("amount", 0, "Expence amount.")

    if err := add.Parse(os.Args[2:]); err != nil {
	fmt.Println("Cannot parse add flags:", err)
	os.Exit(1)
    }

    return *desc, *amount
}
