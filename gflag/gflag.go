// Copyright 2024 da-ket.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file or at
// https://opensource.org/licenses/MIT.

package gflag

// GlobalFlags holds global (PersistentFlags at root command of Cobra) flags.
type GlobalFlags struct {
	SilentDebugMsg bool
}

// G is an accessible variable of type GlobalFlags throughout the package.
var G GlobalFlags
