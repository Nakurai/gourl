package cli

import (
	"fmt"
	"strings"

	"github.com/nakurai/gourl/db"
	"github.com/nakurai/gourl/models"
)

type VarCmd struct{}

// return all the commands that will lead to this execution path
func (c *VarCmd) GetCmds() []string {
	return []string{
		"var",
	}
}

// return all the flags this cmd can handle
func (c *VarCmd) GetFlags() []ValidFlag {
	return []ValidFlag{
		{Key: "data", Labels: []string{"-d", "--data"}},
	}
}

// return all the flags this cmd can handle
func (c *VarCmd) GetHelp() string {
	return `
gourl var list

  List all the variables for the current environment

gourl var add --data key=value

  Create a new variable in the current environment. If the key already exists, the value will be replaced.
	--data, -d: The key/value of the variable, in the format: key=value

gourl var remove --name <name>

  Delete a variable from the current environment.
	--name, -n: The variable's key`
}

// create and send a new http request based on the provided parameters
// we are not expecting any actions here
func (c *VarCmd) Execute(cmd string, actions []string, flags []Flag) (string, error) {
	nbActions := len(actions)
	if nbActions == 0 {
		return fmt.Sprintf("No action provided. You must provide one of the actions below:\n%s\n", c.GetHelp()), nil
	}
	if nbActions > 1 {
		return fmt.Sprintf("Too many actions provided (%d). You must provide one of the actions below:\n%s\n", nbActions, c.GetHelp()), nil
	}

	action := actions[0]
	switch action {
	case "list":
		allVars := ""
		for key, value := range models.CurrentEnv.Variables {
			allVars += fmt.Sprintf("%s: %s\n", key, value)
		}
		return allVars, nil
	case "add":
		newKey := ""
		newValue := ""
		for _, flag := range flags {
			switch flag.Key {
			case "data":
				varParts := strings.Split(flag.Value, "=")
				if len(varParts) != 2 {
					return "", fmt.Errorf("wrong formatting %s. The --data flag must be --data key=value. Use `gourl var` to list all the options", flag.Value)
				}
				newKey = varParts[0]
				newValue = varParts[1]
			default:
				return "", fmt.Errorf("the %s flag is unknown. Use `gourl var` to list all the options", flag.Key)

			}
		}
		if newKey == "" {
			return "", fmt.Errorf("the --data flag is mandatory. Use `gourl var` to list all the options")
		}

		models.CurrentEnv.Variables[newKey] = newValue
		db.Db.Save(&models.CurrentEnv)

		return fmt.Sprintln("done."), nil
	case "remove":
		nameToDelete := ""
		for _, flag := range flags {
			switch flag.Key {
			case "name":
				nameToDelete = flag.Value
			default:
				return "", fmt.Errorf("the %s flag is unknown. Use `gourl var` to list all the options", flag.Key)

			}
		}
		if nameToDelete == "" {
			return "", fmt.Errorf("the --name flag is mandatory. Use `gourl var` to list all the options")
		}

		delete(models.CurrentEnv.Variables, nameToDelete)
		db.Db.Save(&models.CurrentEnv)
		return "Done.", nil

	default:
		return fmt.Sprintf("Invalid action provided (%s). You must provide one of the actions below:\n%s\n", action, c.GetHelp()), nil
	}

}
