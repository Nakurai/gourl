package cli

import "fmt"

// all possible commands must implement this interface
type CmdInterface interface {
	GetCmds() []string
	GetFlags() []ValidFlag
	GetHelp() string
	Execute(cmd string, actions []string, flags []Flag) (string, error)
}

// each command can have a list of flags
type ValidFlag struct {
	Key    string   // the key the Execute function will rely upon, ex: "data"
	Labels []string // possible labels for this flag, ex: "-d", "--data", "--dat" etc
}

// each command's execute function receives a list of flags
type Flag struct {
	Key   string // the key the Execute function will rely upon, ex: "data"
	Value string // the actual value provided by the user
}

// this type is the command's orchestrator
type Cli struct {
	Cmds  map[string]CmdInterface
	Flags map[string]string
	Help  string
}

func NewCli() Cli {
	newCli := Cli{}
	newCli.Cmds = map[string]CmdInterface{}
	newCli.Flags = map[string]string{}
	newCli.Help = "gourl - https://github.com/nakurai/gourl"
	return newCli
}

// this function maps all the possible commands accepted by the program
// it also maps all the flags. this way, we make sure that no two commands are
// trying to use the same keyword, and also that flags are consistent across all commands.
// If d means data somewhere, it means data everywhere.
func (c *Cli) Register(cmds []CmdInterface) error {

	for _, cmd := range cmds {

		// indexing command keywords
		for _, cmdKeyword := range cmd.GetCmds() {
			_, ok := c.Cmds[cmdKeyword]
			if ok {
				return fmt.Errorf("the command keyword %s has already been registered, please change the most recently keyword", cmdKeyword)
			}
			c.Cmds[cmdKeyword] = cmd
		}
		for _, validFlag := range cmd.GetFlags() {
			for _, flagLabel := range validFlag.Labels {
				if flagLabel[0] != '-' {
					return fmt.Errorf("flag %s is ill formatted. All flags must start with a dash '-'", flagLabel)
				}
				existingValue, ok := c.Flags[flagLabel]
				if ok && existingValue != validFlag.Key {
					return fmt.Errorf("a flag with the label '%s' already exist and means '%s', not '%s'. You cannot have two flags with the same label meaning different things. Please change the newest flag's label to something else", flagLabel, existingValue, validFlag.Key)
				}
				c.Flags[flagLabel] = validFlag.Key
			}
		}
		c.Help += "\n\n" + cmd.GetHelp()
	}
	return nil
}

// Parse all the arguments, compare them to the available flags
// It return a list of actions and a map of flags.
// Note: the args parametes EXCLUDES the "cmd". In other words, the first item of the args param is the first actions
// ex: `gourl env add -n test` => args=[add -n test], notice that the env is not passed along
// Note: flags MUST have a value
func (c *Cli) ParseArgs(args []string) ([]string, []Flag, error) {
	actions := []string{}
	flags := []Flag{}

	nbArgs := len(args)
	for argIndex := 0; argIndex < nbArgs; argIndex++ {
		arg := args[argIndex]
		if arg[0] == '-' {
			// this is a flag!
			flagKey, ok := c.Flags[arg]
			if !ok {
				return nil, nil, fmt.Errorf("the flag %s has not been registered by any commands", arg)
			}
			if argIndex == nbArgs-1 {
				return nil, nil, fmt.Errorf("it looks like the flag %s has no value. All flags must be provided a value", arg)
			}
			argIndex += 1
			flags = append(flags, Flag{Key: flagKey, Value: args[argIndex]})

		} else {
			// this argument is an action
			actions = append(actions, arg)
		}
	}
	return actions, flags, nil
}

// this function will look for variable in the flags and replace them with the current env's values.
// it errors out if the variable does not exist in the current environment
// @todo
func (c *Cli) ReplaceVariables(flags []Flag) error {

	return nil
}
