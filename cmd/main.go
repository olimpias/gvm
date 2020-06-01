package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/olimpias/gvm/internal/commands"
	"github.com/olimpias/gvm/internal/filesystem"
)

func main() {
	if len(os.Args) < 2 {
		terminateWithErr(errors.New("show helper"))
	}
	fileManager, err := filesystem.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	var command commands.Command
	//TODO add help
	switch os.Args[1] {
	case "dl":
		ver, err := getVersionArg()
		if err != nil {
			terminateWithErr(err)
		}
		command = commands.NewDLCommand(fileManager, ver)
	case "del":
		ver, err := getVersionArg()
		if err != nil {
			terminateWithErr(err)
		}
		command = commands.NewDelCommand(fileManager, ver)
	case "use":
		ver, err := getVersionArg()
		if err != nil {
			terminateWithErr(err)
		}
		command = commands.NewUseCommand(fileManager, ver)
	case "list":
		command = commands.NewListCommand(fileManager)
	default:
		//err
		fmt.Printf("Unknown command %s \n", os.Args[1])
		os.Exit(2)
	}
	if err := command.Validate(); err != nil {
		terminateWithErr(err)
	}
	if err := command.Apply(); err != nil {
		terminateWithErr(err)
	}
	os.Exit(0)
}

func getVersionArg() (string, error) {
	if len(os.Args) < 3 {
		return "", errors.New("version is required")
	}
	return os.Args[2], nil
}

func terminateWithErr(err error) {
	fmt.Println(err)
	os.Exit(1)
}
