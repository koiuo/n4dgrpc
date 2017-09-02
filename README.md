# n4dgrpc
query namerd via io.l5d.mesh interface

## Usage

```
$ n4dgrpc -a namerd:4321 -t 500ms bind /service/consul /prodnamespace
```

Help is there too
```
$ ./n4dgrpc help                                                                                                  [master] 
n4dgrpc is a CLI application that serves as a client for namerd mesh interface

Usage:
  n4dgrpc [command]

Available Commands:
  bind        bind NAME in NAMESPACE
  help        Help about any command
  resolve     resolve PATH to replica set in NAMESPACE

Flags:
  -a, --address string     address of namerd grpc interface as host:port
        If N4DGRPC_ADDRESS environment variable is set, it is used as default
        value for this flag (default "localhost:4321")
  -h, --help               help for n4dgrpc
  -t, --timeout duration   timeout for command
        Some commands involve multiple calls to namerd. This flag sets global
        time limit (default 1s)

Use "n4dgrpc [command] --help" for more information about a command.
```

## Build & install

### TL;DR

```
$ go get github.com/linkerd/linkerd/mesh/core/src/main/protobuf
$ export l5dproto="$GOPATH/src/github.com/linkerd/linkerd/mesh/core/src/main/protobuf"
$ protoc -I "$l5dproto" "$l5dproto/*.proto" --go_out=plugins=grpc:"$l5dproto"
$ go install github.com/edio/n4dgrpc
```

### Explanation

Tool relies on linkerd protobuf definitions (here'n after just _protos_). Instead of vendoring I decided to fetch those
protos from linkerd repo.

So `go get` will fetch protos from the linkerd repo into `$GOPATH`, but it'll complain that there are no buildable go
files, because there are only protos. We then have to use protobuf compiler to generate go API from protos.

After that tool code can be built in a usual way.

## Known issues

1. Tool doesn't send HTTP/2 `GOAWAY` frame (https://github.com/grpc/grpc-go/issues/460). If you use namerd pre 1.1.3 you'll
see scary stacktrace in its logs on every request.

2. Namerd doesn't allow closing gRPC stream channel, `resolve` command causes exception in namerd.
