package main

import (
	"fmt"
	"os"

	"github.com/nomics-crypto/mig/libmig"
)

func main() {
	err := libmig.Run(os.Args[1:])
	if err != nil {
		fmt.Println(err.Error())
	}
}
