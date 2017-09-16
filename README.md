[![Build Status](https://travis-ci.org/edio/n4dgrpc.svg?branch=master)](https://travis-ci.org/edio/n4dgrpc)

```
      ||  ||  ||\
  .===||  ||  ||===.
      ||  ||  ||    \              /        |             
  .===||  ||  ||===. \     ,.——.  /  |   ,--|   ,--,  ,.-- ,--.   ,--
      ||  ||  ||    \      |   |  ---|  |   |  |   |  |    |   | |
  .===||  ||  ||===. \     '   '     '   '--'   '--|  |    |--'   '--
   \  ||  ||  || |  \                             /        |
    \__\   \   \ |___\                             
        \   \   \|
```

command-line utility to query [namerd](https://github.com/linkerd/linkerd) via
`io.l5d.mesh` interface.

## What is it for?

Think of it as _curl_ for `io.l5d.mesh`.

_namerd_ can expose `io.l5d.thriftNameInterpreter`, `io.l5d.mesh` and
`io.l5d.httpController` interface. It is easy (and often convenient) to check name
resolution in _namerd_ via its `io.l5d.httpController` with _curl_

```
$ curl 'localhost:4380/api/1/resolve/galaxyquest?path=/svc/myservice'
```

However, if `io.l5d.httpController` is not enabled, one has to go to admin ui to
check, whether name can be resolved.

_n4dgrpc_ sends resolve or bind requests to _namerd_  via `io.l5d.mesh`
interface.

It also can be used as a lightweight yet meaningful (as opposed to dumb tcp
handshake) health-check for _namerd_.

## Usage

```
$ n4dgrpc -a localhost:4321 resolve /$/inet/github.com/80 /galaxyquest
192.30.253.113:80
192.30.253.112:80
```

Help is there too
```
$ ./n4dgrpc help
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

`go` is required to build _n4dgrpc_. `GOPATH` environment variable must be set.

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
