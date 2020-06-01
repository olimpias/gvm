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
		terminateWithErr(errors.New("No command is set for 'gvm' \n Run 'gvm help' for usage"))
	}
	fileManager, err := filesystem.New()
	if err != nil {
		terminateWithErr(err)
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
	case "help":
		helper()
	default:
		terminateWithErr(fmt.Errorf("Err: Unknown %s command for 'gvm' \n Run 'gvm help' for usage", os.Args[1]))
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
	os.Stderr.WriteString(fmt.Sprintf("Err: %s  \n", err))
	os.Exit(1)
}

func helper() {
	os.Stdout.WriteString("gvm is a go version controller\n")
	os.Stdout.WriteString("Commands:\n")
	os.Stdout.WriteString("list  list the possible downloaded versions that ready to use.\n")
	os.Stdout.WriteString("dl    downloads the version that you specify to your machine.\n")
	os.Stdout.WriteString("use   uses the version that specify as an input. It has to be downloaded first using dl command.\n")
	os.Stdout.WriteString("del   deletes the version that you specify as an input\n")
	os.Exit(0)
}
