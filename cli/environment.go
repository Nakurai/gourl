package cli

import (
	"fmt"

	"github.com/nakurai/gourl/db"
	"github.com/nakurai/gourl/models"
)

type EnvCmd struct{}

// return all the commands that will lead to this execution path
func (c *EnvCmd) GetCmds() []string {
	return []string{
		"env",
	}
}

// return all the flags this cmd can handle
func (c *EnvCmd) GetFlags() []ValidFlag {
	return []ValidFlag{
		{Key: "name", Labels: []string{"-n", "--name"}},
		{Key: "description", Labels: []string{"-desc", "--description"}},
		{Key: "copy", Labels: []string{"-c", "--copy"}},
	}
}

// return all the flags this cmd can handle
func (c *EnvCmd) GetHelp() string {
	return `
gourl env list

  List all the available environments

gourl env add --name <name> [--copy <name>] [--description <your description>]

  Create a new environment. If the copy flag is used, all variables (and their values) of this other environment will be copied over.
	--name, -n: an arbitrary string to name your environment.
	--copy, -c: an existing environment name.
	--description, -desc: a description of the environment.

gourl env remove --name <name>

  Delete an environment. It will also delete all the variables linked to this environment, and their values. If the deleted environment is the currently loaded one, then the default environment will be loaded automatically.
	--name, -n: an arbitrary string to name your environment.

gourl env load --name <name>

  Will load the environment and its variables for all the following requests.
	--name, -n: an existing environment name.`

}

// create and send a new http request based on the provided parameters
// we are not expecting any actions here
func (c *EnvCmd) Execute(cmd string, actions []string, flags []Flag) (string, error) {
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
		envs := []models.Environment{}
		res := db.Db.Find(&envs)
		if res.Error != nil {
			return "", res.Error
		}
		allEnvs := ""
		for _, env := range envs {
			allEnvs += fmt.Sprintf("%v\n", env)
		}
		return allEnvs, nil
	case "add":
		newName := ""
		newDesc := ""
		copyFrom := ""
		for _, flag := range flags {
			switch flag.Key {
			case "name":
				newName = flag.Value
			case "copy":
				copyFrom = flag.Value
			case "description":
				newDesc = flag.Value
			default:
				return "", fmt.Errorf("the %s flag is unknown. Use `gourl env` to list all the options", flag.Key)

			}
		}
		if newName == "" {
			return "", fmt.Errorf("the --name flag is mandatory. Use `gourl env` to list all the options")
		}

		existingEnv, err := models.GetEnv(newName)
		if err != nil {
			return "", err
		}
		if existingEnv != nil {
			return "", fmt.Errorf("an environment named %s already exists", newName)
		}

		envVars := map[string]string{}
		if copyFrom != "" {
			envToCopy, err := models.GetEnv(copyFrom)
			if err != nil {
				return "", err
			}
			if envToCopy == nil {
				return "", fmt.Errorf("no environment named %s exists, impossible to copy over", copyFrom)
			}
			envVars = envToCopy.Variables
		}

		newEnv := models.Environment{
			Name:      newName,
			Variables: envVars,
			Description: newDesc,
			Current:   false,
		}
		err = models.CreateEnv(&newEnv)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("done. use `gourl env load --name %s` to activate the new environment", newName), nil
	case "remove":
		nameToDelete := ""
		for _, flag := range flags {
			switch flag.Key {
			case "name":
				nameToDelete = flag.Value
			default:
				return "", fmt.Errorf("the %s flag is unknown. Use `gourl env` to list all the options", flag.Key)

			}
		}
		if nameToDelete == "" {
			return "", fmt.Errorf("the --name flag is mandatory. Use `gourl env` to list all the options")
		}
		if nameToDelete == "default" {
			return "", fmt.Errorf("the default env cannot be deleted")
		}

		existingEnv, err := models.GetEnv(nameToDelete)
		if err != nil {
			return "", err
		}
		if existingEnv == nil {
			return "", fmt.Errorf("no environment named %s exists", nameToDelete)
		}

		db.Db.Delete(existingEnv)

		wasExistingEnv := false
		if models.CurrentEnv.Name == nameToDelete{
			wasExistingEnv = true
			models.LoadEnv("default")
		}

		res := fmt.Sprintf("env %s deleted.", nameToDelete)
		if wasExistingEnv{
			res += "Default environment loaded."
		}

		return res, nil
	case "load":
		nameToLoad := ""
		for _, flag := range flags {
			switch flag.Key {
			case "name":
				nameToLoad = flag.Value
			default:
				return "", fmt.Errorf("the %s flag is unknown. Use `gourl env` to list all the options", flag.Key)

			}
		}
		if nameToLoad == "" {
			return "", fmt.Errorf("the --name flag is mandatory. Use `gourl env` to list all the options")
		}

		existingEnv, err := models.GetEnv(nameToLoad)
		if err != nil {
			return "", err
		}
		if existingEnv == nil {
			return "", fmt.Errorf("no environment named %s exists", nameToLoad)
		}

		models.LoadEnv(nameToLoad)

		return fmt.Sprintf("%s loaded", nameToLoad), nil

	default:
		return fmt.Sprintf("Invalid action provided (%s). You must provide one of the actions below:\n%s\n", action, c.GetHelp()), nil
	}

}
