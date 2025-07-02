# Gourl

Gourl (pronounced "goorl") is a command line tool that mixes cURL and tools like postman or insomnia. It allows you to send http requets easily, but also to save and replay them if needed.

## Installation
Download the right executables and make sure it is accessible via your path.

On the first execution, a new folder `.gourl` will be created in the home directory of the user.

## How to use

Commands follow the pattern: `gourl CMD [ACTIONS] [FLAGS]`

### Get queries
The simplest use of gourl is:

`gourl get --url https://jsonplaceholder.typicode.com/posts/1/comments`

The result will be displayed in the console. The output format is:
```
Status

body:
<the body>

headers:
<the headers>
```

To send get parameters, you can use the `--data` flag:

`gourl get --url https://jsonplaceholder.typicode.com/posts/1/comments --data id=1`

### Post queries
Similarly, if you need to send data use the same tag:

`gourl post --url https://jsonplaceholder.typicode.com/posts/1/comments --data id=id1 --data test="my test"`

By default, data will be send as form encoded. To send as JSON, use the `--json` flag:

`gourl post --url https://jsonplaceholder.typicode.com/posts/1/comments --data id=id1 --data test="my test" --json true`

### Headers
For any of those commands, you can specify any header you want by using the flag `--header`:

`gourl post --url https://jsonplaceholder.typicode.com/posts/1/comments --data id=id1 --data test="my test" --json true --header "Authorization=Bearer your-secret-token" --header foo=bar`

### Saving your queries
Retyping the same query everytime is super tedious. If you reuse the same query a lot, you can save the query and retrieve it later. To do that, use the `--save` flag. This flag takes a slash separated name. Each part of the name will be use as a "category", so it is easier to group your queries.
Here is an example:

`gourl post --url [...] --name demo/test/post_message`

In this example, you can think of `demo` and `test` as groups for you queries. For example, you could have another query like:

`gourl post --url [...] --name demo/test/get_message`

This way, it is easy to organize and navigate your saved queries.

## Tip
- A lot of flags have a short form. `-u` for `--url`, `-d` for `--data`, etc. All the forms can be found in via the `help` command.

- If you want to download a file and store its content locally, on Linux you can use the following command: `gourl get --url https://<url-to-file-here> | awk '/^body:/{flag=1; next} /^headers:/{flag=0} flag' > ./filename.txt`

## Roadmap

- [x] Use Interfaces for easily adding more commands
- [x] Support basic ways of sending HTTP requests
- [x] On start, create the sqlite database
- [x] Allow user to save their query
- [x] On start, load the tree of queries in memory
- [x] Add `gourl help` command
- [ ] On start, load the tree of queries in memory with only the first letter of each collection
- [x] Add `gourl list` command to list all queries
- [ ] Add `--depth` flag to the `list` command
- [ ] Add `gourl load --name` command to execute a saved query
- [ ] Add `gourl env list|add|remove` command to create different execution environment
- [ ] Add `gourl var list|add|remove` command to create variables usable in flags
- [ ] Extend vaiables in flags in request logic

## License

Released under the [MIT License](/LICENSE.txt)


