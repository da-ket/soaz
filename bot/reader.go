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
	"strings"

	"github.com/chromedp/chromedp"
)

func ReadPage(keywords []string) (string, error) {
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	var text string
	err := chromedp.Run(ctx,
		chromedp.Navigate(fmt.Sprintf("https://search.naver.com/search.naver?where=nexearch&sm=top_hty&fbm=0&ie=utf8&query=%s", strings.Join(keywords, "+"))),
		chromedp.Text(`div#wrap > div#header_wrap`, &text),
	)
	if err != nil {
		return "", err
	}
	return text, nil
}
