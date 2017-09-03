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

type BindConfig struct {
	FailOnNeg bool
	Root      *mesh.Path
	Name      *mesh.Path
}

var (
	bindConfig = BindConfig{}
)

var bindCmd = &cobra.Command{
	Use:   "bind NAME [NAMESPACE]",
	Short: "bind NAME in NAMESPACE",
	Long: `bind NAME in NAMESPACE

By default command exits with zero even if binding is negative. See options to
change this behavior.
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
			bindConfig.Name = name
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
			bindConfig.Root = root
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		paths, err := client.Bind(bindConfig.Root, bindConfig.Name)
		if err != nil {
			switch err.(type) {
			case *client.ErrNegBinding:
				if bindConfig.FailOnNeg {
					Exit(ExitBindingError, "%v", err)
				}
				break
			default:
				Exit(ExitUnexpectedError, "%v", err)
			}
		}

		for _, path := range paths {
			fmt.Println(convertions.PathToStr(path))
		}
		return
	},
}

func init() {
	N4dgrpc.AddCommand(bindCmd)
	bindCmd.Flags().BoolVarP(&bindConfig.FailOnNeg, "fail", "f", false, "fail if binding is negative")
}
