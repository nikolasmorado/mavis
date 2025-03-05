# Mavis, a resource fetching utility
Fetching with saving requests

Named after my cat mavis who likes to play fetch

```
Usage:
  mavis [method] [url] [flags]
  mavis [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  stash       interact with the stash

Flags:
  -b, --cookie strings   Cookie to send with the request, formatted as Cookie=Value
  -d, --data string      Data to pass along with the request
  -H, --header strings   Headers to be sent with the query formatted as Name:Value 
  -h, --help             help for mavis
  -n, --name string      Name to stash the request as, can use slashes to denote directories
      --stash            Indicates the request should be stashed

Use "mavis [command] --help" for more information about a command.
```

### Todo:
- ~~Stashing, configurable with hcl or toml files or something else, havent decided but dont want json it smells bad~~
- ~~Listing stashed requests~~
- Chaining requests
- Optionally supply data with files
- Support output to files instead of only stdout
- Add picture or video of mavis playing fetch to the readme or something so people believe me
- Make cookies work
- ~~Running stashed requests~~
- Add validation for re running requests
- Prettier output when printing to cli not files
- Websockts
- Supporting other formats not just toml (like hugo)
- Deleting stashed requests
- Stash list specify the project
- Better docs!
- Finish the setting to apply jq to data
