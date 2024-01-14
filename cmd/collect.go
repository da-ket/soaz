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

func dedupKeywords(keywords []string) []string {
	m := make(map[string]string)
	for _, k := range keywords {
		m[k] = ""
	}
	dedupKeywords := make([]string, 0)
	for key, _ := range m {
		dedupKeywords = append(dedupKeywords, key)
	}
	return dedupKeywords
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

			// Check if it is supported platform type or not.
			platformType := bot.StringToPlatformType(collectCmd.platform)
			if platformType == bot.Unsupported {
				collectCmd.err = fmt.Errorf("the platform %s is not supported", collectCmd.platform)
				return
			}

			collectCmd.keywords = dedupKeywords(collectCmd.keywords)

			collectCmd.err = nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if collectCmd.err != nil {
				fmt.Printf("Error: %v\n", collectCmd.err)
				fmt.Println(cmd.UsageString())
				return
			}

			platformType := bot.StringToPlatformType(collectCmd.platform)
			switch platformType {
			case bot.NaverBlog:
				fmt.Println(bot.ReadNaverBlogs(collectCmd.keywords))
			default:
				// Unreachable code.
				// During the pre-run phase, the error-handling for unsupported platform types has been completed.
				// However, it was added for the program's safety.
				panic(fmt.Errorf("the platform (%s) is not supported", collectCmd.platform))
			}
		},
	}
	RootCmd.AddCommand(collectCmd.command)

	collectCmd.command.Flags().StringVarP(&collectCmd.platform, "platform", "p", "", fmt.Sprintf("social-media or search-engine to search keywords from (choose one of from: %s)", bot.SupportedPlatformTypes()))
	collectCmd.command.Flags().StringSliceVarP(&collectCmd.keywords, "keywords", "k", []string{}, "set the keywords to research in deep, it would be your brand or product names separated by a comma (e.g. '--keywords=cocacola,pepsi')")
	collectCmd.command.MarkFlagRequired("platform")
	collectCmd.command.MarkFlagRequired("keywords")
}
