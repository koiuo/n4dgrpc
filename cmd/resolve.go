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
	"github.com/edio/n4dgrpc/client"
	"github.com/edio/n4dgrpc/convertions"
	mesh "github.com/linkerd/linkerd/mesh/core/src/main/protobuf"
	"github.com/spf13/cobra"
)

type ResolveConfig struct {
	FailOnNegBinding   bool
	FailOnEmptyReplica bool
	Root               *mesh.Path
	Name               *mesh.Path
}

var (
	resolveConfig = ResolveConfig{}
)

var resolveCmd = &cobra.Command{
	Use:   "resolve PATH [NAMESPACE]",
	Short: "resolve PATH to replica set in NAMESPACE",
	Long: `resolve PATH to replica set in NAMESPACE

This command invokes 2 operations on namerd: binding and then resolving the
bound name.

See https://linkerd.io/advanced/routing/ for details.

By default command exits with zero even if binding is negative or resolved
replica set is empty. See options to change this behavior.
`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.RangeArgs(1, 2)(cmd, args); err != nil {
			return err
		}

		{
			var name *mesh.Path
			var err error
			if name, err = convertions.StrToPath(args[0]); err != nil {
				return fmt.Errorf("NAME: %v", err)
			}
			resolveConfig.Name = name
		}

		{
			rootStr := DefaultRoot
			if len(args) == 2 {
				rootStr = args[1]
			}

			var root *mesh.Path
			var err error
			if root, err = convertions.StrToPath(rootStr); err != nil {
				return fmt.Errorf("ROOT: %v", err)
			}
			resolveConfig.Root = root
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// TODO error handling
		endpoints, err := client.Resolve(resolveConfig.Root, resolveConfig.Name)
		if err != nil {
			switch err.(type) {
			case *client.ErrNegBinding:
				if resolveConfig.FailOnNegBinding {
					Exit(ExitBindingError, "%v", err)
				}
				break
			default:
				Exit(ExitUnexpectedError, "%v", err)
			}
		}

		if len(endpoints) == 0 && resolveConfig.FailOnEmptyReplica {
			Exit(ExitEmptyReplica, "%v", err)
		}

		for _, endpoint := range endpoints {
			str, _ := convertions.EndpointToStr(endpoint)
			fmt.Println(str)
		}

		return
	},
}

func init() {
	N4dgrpc.AddCommand(resolveCmd)
	resolveCmd.Flags().BoolVarP(&resolveConfig.FailOnNegBinding, "fail-neg", "f", false, "fail if binding is negative")
	resolveCmd.Flags().BoolVarP(&resolveConfig.FailOnEmptyReplica, "fail-empty", "F", false, "fail if replica set is empty")
}
