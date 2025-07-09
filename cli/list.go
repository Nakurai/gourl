package cli

import (
	"github.com/nakurai/gourl/models"
)

type ListCmd struct{}

// return all the commands that will lead to this execution path
func (c *ListCmd) GetCmds() []string {
	return []string{
		"list",
	}
}

// return all the flags this cmd can handle
func (c *ListCmd) GetFlags() []ValidFlag {
	return []ValidFlag{}
}

// return all the flags this cmd can handle
func (c *ListCmd) GetHelp() string {
	return `
gourl list

  List all the queries you have saved.`
}

// create and send a new http request based on the provided parameters
// we are not expecting any actions here
func (c *ListCmd) Execute(cmd string, actions []string, flags []Flag) (string, error) {
	return models.QueryTree.Print(-1, ""), nil
}
