package main

import (
	"fmt"
	"testing"
)

func TestSwitch(t *testing.T) {
    number := 2

    switch {
    case number == 2:
	fmt.Println("number is equal to 2")
    case number % 2 == 0:
	fmt.Println("number is even")
    }
}
