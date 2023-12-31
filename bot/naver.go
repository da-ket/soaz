// Copyright 2023 da-ket.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file or at
// https://opensource.org/licenses/MIT.

package bot

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/chromedp/chromedp"
)

// blog is a struct representing information about a blog.
type blog struct {
	title   string
	link    string
	content string
}

func ReadNaverBlogs(keywords []string) (string, error) {
	const topN = 15

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	// TODO: Retrieve data within the last 3 months
	searchURL := fmt.Sprintf("https://search.naver.com/search.naver?where=blog&query=%s", strings.Join(keywords, "+"))
	err := chromedp.Run(ctx,
		chromedp.Navigate(searchURL),
		chromedp.WaitVisible("#main_pack"),
	)
	if err != nil {
		return "", err
	}

	// Parallel blog data scraping by goroutines.
	ch := make(chan blog, topN)
	for i := 1; i <= topN; i++ {
		go func(i int) {
			b := blog{}
			attrs := make(map[string]string, 0)
			err := chromedp.Run(ctx,
				chromedp.Text(fmt.Sprintf("li#sp_blog_%d div.title_area", i), &b.title),
				chromedp.Attributes(fmt.Sprintf("li#sp_blog_%d div.title_area a", i), &attrs),
			)
			if err != nil {
				panic(err)
			}
			b.link = attrs["href"]

			subCtx, subCancel := chromedp.NewContext(
				context.Background(),
				chromedp.WithDebugf(log.Printf),
			)
			defer subCancel()

			err = chromedp.Run(subCtx,
				chromedp.Navigate(b.link),
				chromedp.WaitVisible("iframe#mainFrame"),
				chromedp.Attributes("iframe#mainFrame", &attrs),
			)
			if err != nil {
				panic(err)
			}

			// We don't want to parse the annoying iframe.
			// Navigate to source page of iframe directly.
			b.link = fmt.Sprintf("http://blog.naver.com%s", attrs["src"])
			err = chromedp.Run(subCtx,
				chromedp.Navigate(b.link),
				chromedp.WaitVisible("div.se-main-container"),
				chromedp.Text(`div.se-main-container`, &b.content),
			)
			if err != nil {
				panic(err)
			}

			// Concatenate all words with a single whitespace.
			removeDuplicateSpaces := regexp.MustCompile(`\s+`)
			b.content = removeDuplicateSpaces.ReplaceAllString(b.content, " ")

			ch <- b
		}(i)
	}

	// Create a slice to store information about the most relevant N blogs in recent 3 months.
	blogs := make([]blog, topN)
	// Wait goroutines.
	for i := 0; i < topN; i++ {
		b := <-ch
		blogs[i] = b
	}

	var resultBuilder strings.Builder
	for i, b := range blogs {
		resultBuilder.WriteString(fmt.Sprintf("[blog info no.%d]\nTitle: %s\nLink: %s\nContent: %s\n\n", i+1, b.title, b.link, b.content))
	}
	result := resultBuilder.String()

	return result, nil
}
