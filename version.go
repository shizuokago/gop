package gop

import (
	"fmt"
	"strings"
)

type Version struct {
	src      string
	version  int
	revision int
	minor    int
	pr       string
	build    string
}

func NewVersion(src string) *Version {
	var ver Version
	ver.version = -1
	ver.revision = -1
	ver.minor = -1
	ver.parse(src)
	return &ver
}

func (v *Version) parse(src string) {

	v.src = src
	remain := ""

	ver := -1
	idx := strings.Index(src, ".")
	if idx == -1 {
		ver = parseInt(v.src)
		remain = "-1.-1"
	} else {
		ver = parseInt(src[:idx])
		remain = src[idx+1:]
	}

	r := -1
	idx = strings.Index(remain, ".")
	if idx == -1 {
		r = parseInt(remain)
		remain = "-1"
	} else {
		r = parseInt(remain[:idx])
		remain = remain[idx+1:]
	}

	idx = strings.Index(remain, "-")
	build := false
	if idx == -1 {
		idx = strings.Index(remain, "+")
		build = true
	}

	m := -1
	if idx == -1 {
		m = parseInt(remain)
		remain = ""
	} else {
		m = parseInt(remain[:idx])
		remain = remain[idx+1:]
	}

	pr := remain
	b := ""
	if build {
		pr = ""
		b = pr
	} else {
		idx = strings.Index(remain, "+")
		if idx != -1 {
			b = pr[:idx]
			pr = pr[idx:]
		}
	}

	//fmt.Println(v.src, ver, r, m, pr, b)

	v.version = ver
	v.revision = r
	v.minor = m

	v.pr = pr
	v.build = b
	return
}

func (v *Version) String() string {
	return fmt.Sprintf("%s", v.src)
}

// 比較
func (v1 *Version) Compare(v2 *Version) int {
	return v1.compare(v2, true)
}

func (v1 *Version) compare(v2 *Version, footer bool) int {

	if v1.version > v2.version {
		return 1
	} else if v1.version < v2.version {
		return -1
	}

	if v1.revision > v2.revision {
		return 1
	} else if v1.revision < v2.revision {
		return -1
	}

	if v1.minor > v2.minor {
		return 1
	} else if v1.minor < v2.minor {
		return -1
	}

	if footer {

		//TODO 存在しない場合
		pIdx := v1.compareString(v1.pr, v2.pr)
		if pIdx != 0 {
			return pIdx
		}

		bIdx := v1.compareString(v1.build, v2.build)
		if bIdx != 0 {
			return bIdx
		}
	}

	return 0
}

// PR、Buildに関しては空文字を大きいと判定
func (v *Version) compareString(str1, str2 string) int {

	if str1 == "" && str2 == "" {
		return 0
	}

	if str1 == "" || str1 > str2 {
		return 1
	} else if str2 == "" || str1 < str2 {
		return -1
	}
	return 0
}

func (v1 *Version) EqCore(v2 *Version) bool {
	return v1.compare(v2, false) == 0
}

func (v1 *Version) Eq(v2 *Version) bool {
	return v1.Compare(v2) == 0
}

func (v1 *Version) Lt(v2 *Version) bool {
	return v1.Compare(v2) < 0
}

func (v1 *Version) Le(v2 *Version) bool {
	return v1.Compare(v2) <= 0
}

func (v1 *Version) Gt(v2 *Version) bool {
	return v1.Compare(v2) > 0
}

func (v1 *Version) Ge(v2 *Version) bool {
	return v1.Compare(v2) >= 0
}

// エラーが怒ってないか？
func (v *Version) IsError() bool {
	if v.version < 0 || v.revision < 0 || v.minor < 0 {
		return true
	}
	return false
}

// 指定バージョンのバージョンかを判定
func (v *Version) IsVersion(val int) bool {
	return v.version == val
}

// PR,Buildが存在しないものをリリース状態とする
func (v *Version) IsRelease() bool {
	if v.IsError() {
		return false
	}
	return v.pr == "" && v.build == ""
}
