package cmd

import (
"github.com/spf13/cobra"
"fmt"
)

var version string

func init() {
	N4dgrpc.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version and exit",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(version)
	},
}
