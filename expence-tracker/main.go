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
    expencesFilepath= "expences.json"
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

    writeExpencesToFile(&expences, expencesFile)

    fmt.Printf("# Expence added successfully (ID: %d)\n", expenceId)   
}

func proceedList() {
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
    expencesFile, expences := getExpences()
    defer expencesFile.Close()

    id := handleDeleteFlags()
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
	writeExpencesToFile(&expences, expencesFile)
    }

    fmt.Println("# Expences deleted successfully")   
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

func handleDeleteFlags() int {
    del := flag.NewFlagSet("del", flag.ExitOnError)

    id := del.Int("id", -1, "Expence id")
    // TODO: add --all flag

    if err := del.Parse(os.Args[2:]); err != nil {
	fmt.Println("Cannot parse delete flags:", err)
	os.Exit(1)
    }

    return *id
}

func getExpences() (*os.File, []Expence) {
    // open file
    expencesFile, err := os.OpenFile(expencesFilepath, os.O_RDWR|os.O_CREATE, 0660)   
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

    // TODO: turn this to a map
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

func writeExpencesToFile(expences *[]Expence, file *os.File) {
    // write updated content back
    json, err := json.Marshal(expences)
    if err != nil {
	fmt.Println("Error marshaling file:", err)
	file.Close()
	os.Exit(1)
    }

    // clean file
    if err = file.Truncate(0); err != nil {
	fmt.Println("Error truncating file:", err)
	file.Close()
	os.Exit(1)
    }
    if _, err = file.Seek(0,0); err != nil {
	fmt.Println("Error seeking file:", err)
	file.Close()
	os.Exit(1)
    }

    // write updated content to a file
    if _, err = file.Write(json); err != nil {
	fmt.Println("Error writing to file:", err)
	file.Close()
	os.Exit(1)
    }
}
