// Copyright 2023 da-ket.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file or at
// https://opensource.org/licenses/MIT.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// This is the collect command.
var collectCmd struct {
	platform string
	command  *cobra.Command
}

func init() {
	collectCmd.command = &cobra.Command{
		Use:   "collect",
		Short: "Collect the meaningful data on the web",
		Long:  "Collect the meaningful data on the web",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(collectCmd.platform)
		},
	}
	RootCmd.AddCommand(collectCmd.command)
	collectCmd.command.Flags().StringVarP(&collectCmd.platform, "platform", "p", "", "social-media or search-engine to search keywords from")
}
