package main

import "fmt"
import "os"

const (
	dirName = "migrations"
)

// ErrUsage is returned if mig is used improperly
var ErrUsage = fmt.Errorf(`usage:
	mig <command> [arguments]

commands:
	init
	help`)

// Run is the entrypoint for the executable. It takes and arguments list
// and returns an error.
func Run(args []string) error {
	if len(args) == 0 {
		return ErrUsage
	}
	switch args[0] {
	case "init":
		return Init()
	}
	return ErrUsage
}

// Init is called via `mig init` and initializes a project for mig.
func Init() error {
	err := os.Mkdir(dirName, 0755)
	if err == nil {
		fmt.Println("mig initialized")
	}
	return err
}
