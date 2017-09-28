package main

import (
	"fmt"
	"os"
)

func main() {
	err := Run(os.Args[1:])
	if err != nil {
		fmt.Println(err.Error())
	}
}
