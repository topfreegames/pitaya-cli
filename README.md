pitaya-cli
==================

A REPL cli client made in go for pitaya.

### installing

```
go get -u github.com/topfreegames/pitaya-cli
```

### usage
```
$ pitaya-cli

Pitaya REPL Client
>>> help

Commands:
  clear           clear the screen
  connect         connects to pitaya
  disconnect      disconnects from pitaya server
  exit            exit the program
  help            display help
  notify          makes a notify to pitaya server
  request         makes a request to pitaya server
```

For connecting to a protobuffer server specify the documentation route with the -docs argument:

```
pitaya-cli -docs connector.docsHandler.docs
```

Protobuffer servers must implement remotes for protobuf descriptors and auto documentation.
