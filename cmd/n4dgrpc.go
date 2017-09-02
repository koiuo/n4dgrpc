// Copyright © 2017 Dmytro Kostiuchenko edio@archlinux.us
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
	"time"
	"github.com/edio/n4dgrpc/client"
)

// exit codes
type ExitCode uint8

const (
	// Error during parsing input data, validation, etc.
	ExitBindingError = ExitCode(2)
	ExitEmptyReplica = ExitCode(3)
	ExitUnexpectedError = ExitCode(127)
)

// shared const
const (
	DefaultRoot = "/default"
)

type N4dgrpcConfig struct {
	Timeout time.Duration
	Address string
}

var (
	n4dgrpcConfig = N4dgrpcConfig{}
)

var N4dgrpc = &cobra.Command{
	Use:   "n4dgrpc",
	Short: "grpc client for namerd",
	Long: "n4dgrpc is a CLI application that serves as a client for namerd mesh interface",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		client.OpTimeout = n4dgrpcConfig.Timeout
		if err := client.Connect(n4dgrpcConfig.Address); err != nil {
			Exit(ExitUnexpectedError, "Connection error: %v", err)
		}
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		client.Close()
	},
}

func Execute() {
	if err := N4dgrpc.Execute(); err != nil {
		Exit(ExitUnexpectedError, "Error: %v", err)
	}
}

func init() {
	N4dgrpc.PersistentFlags().StringVarP(&n4dgrpcConfig.Address, "address", "a", os.Getenv("N4DGRPC_ADDRESS"),
			`address of namerd grpc interface as host:port
	If N4DGRPC_ADDRESS environment variable is set, it is used as default
	value for this flag`)
	N4dgrpc.PersistentFlags().DurationVarP(&n4dgrpcConfig.Timeout, "timeout", "t", time.Second,
			`timeout for command
	Some commands involve multiple calls to namerd. This flag sets global
	time limit`)
}

func Exit(code ExitCode, msg string, args ...interface{}) {
	if len(args) > 0 {
		fmt.Fprintf(os.Stderr, msg + "\n", args)
	} else {
		fmt.Fprintln(os.Stderr, msg)
	}
	os.Exit(int(code))
}
