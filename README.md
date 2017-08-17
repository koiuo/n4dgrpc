# n4dgrpc
query namerd via io.l5d.mesh interface

## Usage

```
$ n4dgrpc -a namerd:4321 -t 500ms bind /service/consul /prodnamespace
```

Help is there too
```
$ n4dgrpc --help
  n4dgrpc is a CLI application for invoking various operations on
           namerd via its mesh interface.
  
  Usage:
    n4dgrpc [command]
  
  Available Commands:
    bind        bind name in namespace
    help        Help about any command
    resolve     resolve path to replica set in namespace
  
  Flags:
    -a, --address string     address of namerd grpc interface as host:port.
                  If N4DGRPC_ADDRESS environment variable is specified, its value is used by default.
    -h, --help               help for n4dgrpc
    -t, --timeout duration   wait no longer than specified time for command to complete.
                  Some commands involve multiple calls to namerd. This flag sets global time limit. (default 1s)
  
  Use "n4dgrpc [command] --help" for more information about a command.
```
