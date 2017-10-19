package main

import (
	"log"
	"os"

	"github.com/nomics-crypto/mig/libmig"
)

func main() {
	err := libmig.Run(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
}
