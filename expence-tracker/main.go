package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

var (
    availableRecordId int
    MIN_CMD_ARGS = 2 
)

func init () {
    // TODO: lastId
}

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
    Id int
    Description string
    Amount int
}

func proceedAdd() {
    desc, amount := handleAddFlags()

    expencesFile, err := os.OpenFile("expences.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)   
    if err != nil {
	fmt.Printf("Couldn't open expences.json: %s\n", err)
	os.Exit(1)
    }
    defer expencesFile.Close()

    availableRecordId += 1
    record := Expence {
	Id: availableRecordId,
	Description: desc,
	Amount: amount,
    }

    recordJson, err := json.Marshal(record)
    if err != nil {
	fmt.Printf("Couldn't marshal new object: %s\n", err)
	expencesFile.Close()
	os.Exit(1)
    }

    _, err = expencesFile.Write(recordJson)
    if err != nil {
	fmt.Printf("Record not added: %s\n", err)
	expencesFile.Close()
	os.Exit(1)
    }

    fmt.Printf("# Expence added successfully (ID: %d)\n", availableRecordId)   
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
	fmt.Printf("Cannot parse add flags. %s", err)
	os.Exit(1)
    }

    return *desc, *amount
}
