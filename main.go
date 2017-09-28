package main

import (
	"fmt"
	"os"

	"github.com/ngauthier/mig/libmig"
)

func main() {
	err := libmig.Run(os.Args[1:])
	if err != nil {
		fmt.Println(err.Error())
	}
}
