package main

import (
	"fmt"
	"os"
	"qqweq/playground/expence-tracker/cmd"
)

func main() {
    if err := cmd.Execute(); err != nil {
	fmt.Fprint(os.Stderr, err)
	os.Exit(1)
    }
}
