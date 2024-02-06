// Copyright 2024 da-ket.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file or at
// https://opensource.org/licenses/MIT.

package cmd

import (
	"strings"
	"testing"
)

func TestDedupKeywords(t *testing.T) {
	testCases := []struct {
		input  []string
		output int
	}{
		{[]string{"감자", "고구마", "감자"}, 2},
		{[]string{"짜장면", "짬뽕", "짬뽕", "짬뽕"}, 2},
	}

	for _, tc := range testCases {
		out := len(dedupKeywords(tc.input))
		if out != tc.output {
			t.Errorf("`%s`: expected `%d`, got `%d`", strings.Join(tc.input, ","), tc.output, out)
		}
	}
}
