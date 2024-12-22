package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"qqweq/playground/expence-tracker/models"
)

var (
    MIN_CMD_ARGS = 2 
    expencesFilepath= "expences.json"
)

func getExpences() (*os.File, []models.Expence) {
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
    var expences []models.Expence 

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

func writeExpencesToFile(expences *[]models.Expence, file *os.File) {
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
