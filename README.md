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

Protobuffer servers must implement handlers for protobuf descriptors and auto documentation.

A full example of running pitaya-cli with protobuffers:

```
pitaya-cli -docs connector.docsHandler.docs
>>> push connector.playerHandler.matchfound protos.FindMatchPush
>>> connect localhost:30124
>>> request connector.playerHandler.create
>>> request connector.playerHandler.findmatch {"RoomType":"xxxx"}
```

This example can be run with the pitaya-bot protobuffer testing example.
