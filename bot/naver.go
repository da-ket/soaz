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
	"sync"
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
	err     error
}

func ReadNaverBlogs(keywords []string, topN int) (string, error) {
	if gflag.G.SilentDebugMsg {
		log.SetOutput(io.Discard)
	}

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithDebugf(log.Printf),
	)
	go func() {
		time.Sleep(time.Duration(60*(1+(topN/30))) * time.Second)
		cancel()
	}()
	defer cancel()

	scrollingScript := fmt.Sprintf(`
		let nScroll = 0
		const scrollInterval = setInterval(() => {
			window.scrollTo(0, document.body.scrollHeight)
			if (++nScroll === %d) {
				clearInterval(scrollInterval)
			}
		}, %d)
	`, topN/30, 300) // Scroll down times, wait ms per scroll.

	// TODO: Retrieve data within the last 3 months
	searchURL := fmt.Sprintf("https://search.naver.com/search.naver?ssc=tab.blog.all&sm=tab_jum&query=%s", strings.Join(keywords, "+"))
	err := chromedp.Run(ctx,
		chromedp.Navigate(searchURL),
		chromedp.Evaluate(scrollingScript, nil),
		chromedp.Sleep(time.Duration(1+(topN/30))*time.Second),
		chromedp.WaitVisible("#main_pack"),
	)
	if err != nil {
		return "", err
	}

	// Adjusts 'topN' based on the actual number of blogs found.
	var nodes []*cdp.Node
	err = chromedp.Run(ctx,
		chromedp.Nodes("div.title_area > a", &nodes, chromedp.ByQueryAll),
	)
	if err != nil {
		return "", err
	}
	if len(nodes) < topN {
		topN = len(nodes)
	}
	fmt.Printf("%d blog links are fetched. First %d links will be navigated..\n", len(nodes), topN)

	var mu sync.Mutex
	var runningGoRoutines int

	// Parallel blog data scraping by goroutines.
	const numGoroutinesLimit = 20
	ch := make(chan blog, topN)
	for i := 0; i < topN; i++ {
		go func(i int) {
			for {
				mu.Lock()
				if runningGoRoutines < numGoroutinesLimit {
					runningGoRoutines++
					mu.Unlock()
					break
				} else {
					mu.Unlock()
					time.Sleep(100 * time.Millisecond)
				}
			}
			defer func() {
				mu.Lock()
				runningGoRoutines--
				mu.Unlock()
			}()

			b := blog{}
			b.err = chromedp.Run(ctx,
				chromedp.Text([]cdp.NodeID{nodes[i].NodeID}, &b.title, chromedp.ByNodeID),
			)
			if b.err != nil {
				b.err = fmt.Errorf("Parent request meet an error: %w", b.err)
				ch <- b
				return
			}

			b.link = nodes[i].AttributeValue("href")
			if !strings.HasPrefix(b.link, "https://blog.naver.com") {
				b.err = fmt.Errorf("Child request meet an error: Invalid URL")
				ch <- b
				return
			}

			subCtx, subCancel := chromedp.NewContext(
				context.Background(),
				chromedp.WithDebugf(log.Printf),
			)
			go func() {
				time.Sleep(60 * time.Second)
				subCancel()
			}()

			attrs := make(map[string]string)
			b.err = chromedp.Run(subCtx,
				chromedp.Navigate(b.link),
				chromedp.WaitVisible("iframe#mainFrame"),
				chromedp.Attributes("iframe#mainFrame", &attrs),
			)
			if b.err != nil {
				b.err = fmt.Errorf("Child request meet an error: %w", b.err)
				subCancel()
				ch <- b
				return
			}

			// We don't want to parse the annoying iframe.
			// Navigate to source page of iframe directly.
			b.link = fmt.Sprintf("http://blog.naver.com%s", attrs["src"])
			b.err = chromedp.Run(subCtx,
				chromedp.Navigate(b.link),
				chromedp.WaitVisible("div.se-main-container"),
				chromedp.Text(`div.se-main-container`, &b.content),
			)
			if b.err == context.Canceled {
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
				b.err = chromedp.Run(subCtx,
					chromedp.Navigate(b.link),
					chromedp.WaitVisible("#postViewArea"),
					chromedp.Text(`#postViewArea`, &b.content),
				)
				if b.err == context.Canceled {
					// Few blogs have 'div.se_paragraph' rather than '#postViewArea'.
					// TODO(da-ket): clean up this.
					subCtx, subCancel = chromedp.NewContext(
						context.Background(),
						chromedp.WithDebugf(log.Printf),
					)
					go func() {
						time.Sleep(15 * time.Second)
						subCancel()
					}()
					b.err = chromedp.Run(subCtx,
						chromedp.Navigate(b.link),
						chromedp.WaitVisible("div.sect_dsc"),
						chromedp.Text(`div.sect_dsc`, &b.content),
					)
				}
			}
			if b.err != nil {
				b.err = fmt.Errorf("Child request meet an error: %w", b.err)
				subCancel()
				ch <- b
				return
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
	blogs := make([]blog, 0, topN)
	// Wait goroutines.
	for i := 1; i <= topN; i++ {
		b := <-ch
		if b.err != nil {
			fmt.Printf("[%.03d - Err]: %v (URL: %v)\n", i, b.err, b.link)
		} else {
			fmt.Printf("[%.03d - Fin]: %v\n", i, b.link)
			blogs = append(blogs, b)
		}
	}

	var resultBuilder strings.Builder
	for i, b := range blogs {
		resultBuilder.WriteString(fmt.Sprintf("[blog info no.%d]\nTitle: %s\nLink: %s\nContent: %s\n\n", i+1, b.title, b.link, b.content))
	}
	result := resultBuilder.String()

	return result, nil
}
