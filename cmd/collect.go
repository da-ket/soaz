// Copyright 2023 da-ket.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file or at
// https://opensource.org/licenses/MIT.

package cmd

import (
	"fmt"
	"unicode/utf8"

	"github.com/da-ket/soaz/bot"
	"github.com/spf13/cobra"
)

// This is the collect command.
var collectCmd struct {
	platform string
	keywords []string
	err      error
	command  *cobra.Command
}

func init() {
	collectCmd.command = &cobra.Command{
		Use:   "collect",
		Short: "Collect the meaningful data on the web",
		Long:  "Collect the meaningful data on the web",
		PreRun: func(cmd *cobra.Command, args []string) {
			const keywordsNumLimit int = 3
			if len(collectCmd.keywords) > keywordsNumLimit {
				collectCmd.err = fmt.Errorf("the number of keywords is limited to %d, got %d", keywordsNumLimit, len(collectCmd.keywords))
				return
			}

			const keywordLenLimit int = 25
			for _, k := range collectCmd.keywords {
				if utf8.RuneCountInString(k) > keywordLenLimit {
					collectCmd.err = fmt.Errorf("each keyword is limited to %d letters, got %d", keywordLenLimit, utf8.RuneCountInString(k))
					return
				}
			}

			// TODO (da-ket): Platform flag is mandatory.
			// TODO (da-ket): Keywords flag is mandatory.
			// TODO (da-ket): Unsupported platforms should return error.
			// TODO (da-ket): Deduplicate keywords should be handled.

			collectCmd.err = nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if collectCmd.err != nil {
				fmt.Printf("Error: %v\n", collectCmd.err)
				fmt.Println(cmd.UsageString())
				return
			}
			fmt.Println(bot.ReadPage(collectCmd.keywords))
		},
	}
	RootCmd.AddCommand(collectCmd.command)

	collectCmd.command.Flags().StringVarP(&collectCmd.platform, "platform", "p", "", "social-media or search-engine to search keywords from")
	collectCmd.command.Flags().StringSliceVarP(&collectCmd.keywords, "keywords", "k", []string{}, "set the keywords to research in deep, it would be your brand or product names separated by a comma (e.g. '--keywords=cocacola,pepsi')")
}
