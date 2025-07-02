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


// return all the flags this cmd can handle
func (c *RequestCmd) GetHelp() string {
	return `
gourl connect|delete|get|head|options|patch|post|put|trace --url <url> [--data test=test] [--header test=test] [--json true|false] [--save <name>] 

  Send a request to the url provided via the --url flags. List of flags are:
    --url, -u: The URL you want to send the request to. This flag is mandatory. Ex: --url https://example.com
    --data, -d: The data you want to send with the request, either in the body or in the query. The format is key=value. If your value has spaces in it, you must surround the value by double quotes. The data will be send in the body for the following methods: POST, PUT, PATCH. You can use this flag several times. Ex: --data test=test -d test2=test2
    --header, -h: You can specify the request header. Format is like the --data flag: key=value. Ex: Authentication="Bearer XYZ"
    --json, -j: if the value is true, then the body will be formatted as a JSON object and the content-type header will be set to application/json (if no content type header was explictly provided)
    --save, -s: you can provide any name here and the query will be save alongside with all the flags. If a query with the same name already exists, it will let you know and not save it.
`
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
