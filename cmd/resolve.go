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
	"github.com/spf13/cobra"
	"github.com/edio/n4dgrpc/client"
	"fmt"
	"github.com/edio/n4dgrpc/convertions"
)

var resolveCmd = &cobra.Command{
	Use:   "resolve PATH [NAMESPACE]",
	Short: "resolve path to replica set in namespace",
	Long: "resolve path to replica set in namespace",
	Args: bindCmd.Args, // just copy from bind for now. TODO think of proper sharing later
	Run: func(cmd *cobra.Command, args []string) {
		// TODO error handling
		endpoints, err := client.Resolve(bindConfig.Root, bindConfig.Name)
		if err != nil {
			Exit(ExitBindingError, "Failed to resolve: %v", err)
		}
		if len(endpoints) == 0 {
			Exit(ExitBindingError, "No replicas resolved")
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
}
