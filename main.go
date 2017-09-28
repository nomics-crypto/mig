package main

import (
	"fmt"
	"os"
)

func main() {
	err := Run(os.Args[1:])
	if err != nil {
		fmt.Print(err.Error())
	}
}
