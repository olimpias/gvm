package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/olimpias/gvm/cmd"
	"github.com/olimpias/gvm/common"
)

func main() {
	if len(os.Args) < 2 {
		terminateWithErr(errors.New("show helper"))
	}
	fileManager, err := common.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	var command cmd.Command
	switch os.Args[1] {
	case "dl":
		ver, err := getVersionArg()
		if err != nil {
			terminateWithErr(err)
		}
		command = cmd.NewDLCommand(fileManager, ver)
	case "del":
		ver, err := getVersionArg()
		if err != nil {
			terminateWithErr(err)
		}
		command = cmd.NewDelCommand(fileManager, ver)
	case "use":
		ver, err := getVersionArg()
		if err != nil {
			terminateWithErr(err)
		}
		command = cmd.NewUseCommand(fileManager, ver)
	case "list":
		command = cmd.NewListCommand(fileManager)
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
