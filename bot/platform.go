// Copyright 2024 da-ket.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file or at
// https://opensource.org/licenses/MIT.

package bot

type PlatformType int

// Collection of supported platform types.
const (
	Unsupported PlatformType = iota
	NaverBlog
	NaverCafe
)

func (t PlatformType) String() string {
	switch t {
	case NaverBlog:
		return "naverblog"
	case NaverCafe:
		return "navercafe"
	default:
		return "unsupported"
	}
}

func StringToPlatformType(str string) PlatformType {
	switch str {
	case NaverBlog.String():
		return NaverBlog
	// TODO (da-ket): we currently only support naver blogs.
	default:
		return Unsupported
	}
}

func SupportedPlatformTypes() string {
	return NaverBlog.String()
	// TODO: return NaverBlog.String() + "," + NaverCafe.String()
}
