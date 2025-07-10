package cli

import (
	"fmt"

	"github.com/nakurai/gourl/models"
)

type LodaCmd struct{}

// return all the commands that will lead to this execution path
func (c *LodaCmd) GetCmds() []string {
	return []string{
		"load",
	}
}

// return all the flags this cmd can handle
func (c *LodaCmd) GetFlags() []ValidFlag {
	return []ValidFlag{
		{Key: "name", Labels: []string{"-n", "--name"}},
	}
}

// return all the flags this cmd can handle
func (c *LodaCmd) GetHelp() string {
	return `
gourl load --name <name>

  Load and execute the saved query using the current environment's variables if necessary`
}

func (c *LodaCmd) Execute(cmd string, actions []string, flags []Flag) (string, error) {
	nameToLoad := ""
	for _, flag := range flags {
		switch flag.Key {
		case "name":
			nameToLoad = flag.Value
		default:
			return "", fmt.Errorf("the %s flag is unknown. Use `gourl load` to list all the options", flag.Key)

		}
	}
	if nameToLoad == "" {
		return "", fmt.Errorf("the --name flag is mandatory. Use `gourl load` to list all the options")
	}

	query, err := models.GetQuery(nameToLoad)
	if err != nil {
		return "", err
	}
	if query == nil {
		return "", fmt.Errorf("no query named %s exists", nameToLoad)
	}

	return query.Send()
}
