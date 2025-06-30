package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type RequestCmd struct{}

// return all the commands that will lead to this execution path
func (c *RequestCmd) GetCmds() []string {
	return []string{
		"connect",
		"delete",
		"get",
		"head",
		"options",
		"patch",
		"post",
		"put",
		"trace",
	}
}

// return all the flags this cmd can handle
func (c *RequestCmd) GetFlags() []ValidFlag {
	return []ValidFlag{
		{Key: "data", Labels: []string{"-d", "--data"}},
		{Key: "header", Labels: []string{"-h", "--header"}},
		{Key: "json", Labels: []string{"-j", "--json"}},
		{Key: "save", Labels: []string{"-s", "--save"}},
		{Key: "url", Labels: []string{"-u", "--url"}},
	}
}

// create and send a new http request based on the provided parameters
// we are not expecting any actions here
func (c *RequestCmd) Execute(cmd string, actions []string, flags []Flag) (string, error) {
	newQuery := Query{
		Method: strings.ToUpper(cmd),
		Data:   map[string]string{},
		Header: map[string]string{},
	}
	for _, flag := range flags {
		switch flag.Key {
		case "data":
			dataParts := strings.Split(flag.Value, "=")
			dataKey := strings.ToLower(dataParts[0])
			newQuery.Data[dataKey] = strings.Join(dataParts, "=")
		case "header":
			dataParts := strings.Split(flag.Value, "=")
			dataKey := strings.ToLower(dataParts[0])
			newQuery.Header[dataKey] = strings.Join(dataParts, "=")
		case "json":
			newQuery.IsJson = flag.Value == "true"
		case "save":
			newQuery.Name = flag.Value
		case "url":
			newQuery.Url = flag.Value
		default:
			return "", fmt.Errorf("unknown flag %s. Use gourl help for a list of valid flags", flag.Key)

		}
	}

	// if the user already specified the Content-Type header, we do not want to override it
	_, contentTypeSet := newQuery.Header["content-type"]
	if !contentTypeSet {
		if newQuery.IsJson {
			newQuery.Header["content-type"] = "application/json"
		} else {
			if newQuery.Method == "POST" || newQuery.Method == "PATCH" || newQuery.Method == "PUT" {
				newQuery.Header["content-type"] = "application/x-www-form-urlencoded"
			}
		}
	}

	if newQuery.Url == "" {
		return "", fmt.Errorf("no url provided. Please specify a url by using the '--url' flag")
	}

	if newQuery.Name != "" {
		err := newQuery.Save()
		if err != nil {
			if err.Error() == "EXIST-ALREADY" {
				// @todo ask the user if they want to override the existing query
				fmt.Println("A query with this name already exists, not saving. Press enter to continue")
				bufio.NewReader(os.Stdin).ReadBytes('\n')
				fmt.Println("Resuming request")
			} else {
				return "", fmt.Errorf("error while saving the query: %v", err)
			}
		}

	}

	return newQuery.Send()

}
