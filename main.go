package main

import (
	"fmt"
	"os"

	"github.com/nakurai/gourl/cli"
	"github.com/nakurai/gourl/db"
	"github.com/nakurai/gourl/utils"
)

// when compiling the program, use go build -ldflags "-X main.version=X.X.X" to
// update this to the correct value
var version string = ""

func init() {
	if version == "" {
		version = "dev"
	}

	err := utils.CreateDataDir()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	err = db.Init()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	err = db.Db.AutoMigrate(&cli.Query{})
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	cli.BuildQueryTree()
	// cli.QueryTree.Print(-1, "")

}

func main() {
	
	if len(os.Args) == 1 {
		fmt.Printf("gourl v.%s - https://github.com/nakurai/gourl\nUse gourl help for doc\n", version)
		return
	}

	app := cli.NewCli()
	err := app.Register([]cli.CmdInterface{&cli.RequestCmd{}})
	if err != nil {
		fmt.Printf("error while registering cmds: %v\n", err)
		return

	}

	cmd := os.Args[1]

	if cmd == "help"{
		fmt.Println(app.Help)
		return
	}

	_, ok := app.Cmds[cmd]
	if !ok {
		errStr := fmt.Sprintf("error, no registered command has the available cmd: %s.\n", cmd)
		if version == "dev" {
			errStr += "Dev: did you forget to register a new command?"
		}
		fmt.Println(errStr)
		return
	}

	actions, flags, err := app.ParseArgs(os.Args[2:])
	if err != nil {
		fmt.Printf("error parsing the argument(s): %v\n", err)
		return
	}

	err = app.ReplaceVariables(flags)
	if err != nil {
		fmt.Printf("error replacing variables: %v\n", err)
		return
	}
	res, err := app.Cmds[cmd].Execute(cmd, actions, flags)
	if err != nil {
		fmt.Printf("error executing the command: %v\n", err)
		return
	}
	fmt.Println(res)

}
