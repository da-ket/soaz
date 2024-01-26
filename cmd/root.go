// Copyright 2023 da-ket.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file or at
// https://opensource.org/licenses/MIT.

package cmd

import (
	"fmt"
	"os"

	"github.com/da-ket/soaz/gflag"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "soaz",
	Short: "Soaz is a fantastic data collector for digital marketers",
	Long:  "Soaz is a fantastic data collector for digital marketers",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&gflag.G.SilentDebugMsg, "quiet", "q", true, "set the quiet mode to suppress debug message from command line output")
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
