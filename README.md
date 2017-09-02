```
      ||  ||  ||\
  .===||  ||  ||===.
      ||  ||  ||    \              /        |             
  .===||  ||  ||===._\     ,.——.  /  |   ,--|   ,--,  ,.-- ,--.   ,--
      ||  ||  ||    \      |   |  ---|  |   |  |   |  |    |   | |
  .===||  ||  ||===._\     '   '     '   '--'   '--|  |    |--'   '--
   \  ||  ||  || |  \                             /        |
    \__\   \   \ |___\                             
        \   \   \|
```

command-line utility to query [namerd](https://github.com/linkerd/linkerd) via
`io.l5d.mesh` interface

[![Build Status](https://travis-ci.org/edio/n4dgrpc.svg?branch=master)](https://travis-ci.org/edio/n4dgrpc)

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

For building `go` is requird. `GOPATH` environment variable must be set.

### From sources using GNU make

For convenience `Makefile` is provided. Building should not require any
intervention. Make script will fetch fresh copy of gRPC service definitions
from _linkerd_ repo, download protobuf compiler, compile gRPC stubs and finally
build an executable binary as `./n4dgrpc`

```
$ make
```

### Using go get/install

Tool relies on _linkerd_ gRPC service definitions.

Not only manual compilation of gRPC stubs is required, but they also should be
fetched from the _linkerd_ repo. Thus installation with single `go get` command
is not possible.

```
$ # Fetch gRPC service definitions.
$ # This command should fail with 'no buildable sources' error. This is normal
$ go get github.com/linkerd/linkerd/mesh/core/src/main/protobuf
$ export l5dproto="$GOPATH/src/github.com/linkerd/linkerd/mesh/core/src/main/protobuf"
$ ./protoc -I "$l5dproto" "$l5dproto/*.proto" --go_out=plugins=grpc:"$l5dproto"
```

Only after this go get/install will succeed

```
$ go install github.com/edio/n4dgrpc
```

## Known issues

1. Tool doesn't send HTTP/2 `GOAWAY` frame
(https://github.com/grpc/grpc-go/issues/460). Thus every operation will cause
an exception logged for _namerd_ < 1.1.3.

2. _namerd_ doesn't allow closing gRPC stream channel, `resolve` will causes
exception logged in namerd.
