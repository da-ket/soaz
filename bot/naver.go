// Copyright 2023 da-ket.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file or at
// https://opensource.org/licenses/MIT.

package bot

import (
	"context"
	"fmt"
	"io"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"github.com/da-ket/soaz/gflag"
)

// blog is a struct representing information about a blog.
type blog struct {
	title   string
	link    string
	content string
}

func ReadNaverBlogs(keywords []string) (string, error) {
	topN := 15

	if gflag.G.SilentDebugMsg {
		log.SetOutput(io.Discard)
	}

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithDebugf(log.Printf),
	)
	go func() {
		time.Sleep(60 * time.Second)
		cancel()
	}()
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

	// Adjusts 'topN' based on the actual number of blogs found.
	var nodes []*cdp.Node
	err = chromedp.Run(ctx,
		chromedp.Nodes("ul.lst_view > li.bx", &nodes, chromedp.ByQueryAll),
	)
	if err != nil {
		return "", err
	}
	if len(nodes) < topN {
		topN = len(nodes)
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
			go func() {
				time.Sleep(30 * time.Second)
				subCancel()
			}()

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
			if err == context.Canceled {
				// Some blogs have '#postViewArea' rather than 'dev.se-main-container'.
				// TODO(da-ket): clean up this.
				subCtx, subCancel = chromedp.NewContext(
					context.Background(),
					chromedp.WithDebugf(log.Printf),
				)
				go func() {
					time.Sleep(30 * time.Second)
					subCancel()
				}()
				err = chromedp.Run(subCtx,
					chromedp.Navigate(b.link),
					chromedp.WaitVisible("#postViewArea"),
					chromedp.Text(`#postViewArea`, &b.content),
				)
			}
			if err != nil {
				panic(err)
			}

			// Concatenate all words with a single whitespace.
			removeDuplicateSpaces := regexp.MustCompile(`\s+`)
			b.content = removeDuplicateSpaces.ReplaceAllString(b.content, " ")

			ch <- b

			// Cancel the subquery context.
			subCancel()
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
