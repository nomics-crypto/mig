package main

import "fmt"
import "os"

const (
	dirName = "migrations"
)

var UsageError = fmt.Errorf(`Usage:
	mig <command> [arguments]

Commands:
	init
	help
`)

func Run(args []string) error {
	if len(args) == 0 {
		return UsageError
	}
	switch args[0] {
	case "init":
		return Init()
	}
	return UsageError
}

func Init() error {
	return os.Mkdir(dirName, 0755)
}
