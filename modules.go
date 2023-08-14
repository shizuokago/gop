// "github.com/hajimehoshi/ebiten/v2/inpututil"のようなパッケージ名称から、
// 実パスを割り出す用の構造体
// v指定があるディレクトリをすべて候補にする
// 実際にはその可能性を総当たりして、バージョン指定を決定する
package gop

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

type Module struct {
	src          string
	packageNames []*pName

	current *pName
	version *Version
}

func (m *Module) String() string {
	var b strings.Builder
	b.WriteString("Module Name:" + m.src)
	for _, n := range m.packageNames {
		b.WriteString("\n  " + n.String())
	}
	return b.String()
}

func (m1 *Module) Eq(m2 *Module) bool {
	if m1.src == m2.src {
		return true
	}
	return false
}

// 指定されている可能性のあるバージョンとその時のパッケージ名
type pName struct {
	name    string
	version int
}

func (p *pName) String() string {
	return fmt.Sprintf("[%s]:Version:%d", p.name, p.version)
}

func newPackageName(n string, v int) *pName {
	var rtn pName
	rtn.name = n
	rtn.version = v
	return &rtn
}

func createCandidatePackageNames(name string) []*pName {

	names := make([]*pName, 0)
	names = append(names, newPackageName(name, -1))

	ds := strings.Split(name, "/")
	leng := 0

	for _, d := range ds {
		//省略したものを追加
		v := parseTargetVersion(d)
		if v != -1 {
			p := "/" + d
			rename := strings.Replace(name, p, "", 1)
			//バージョンの可能性のあるディレクトリを省略して登録
			names = append(names, newPackageName(rename, v))
		}
		//合計文字列長を追加
		// TODO 位置を算出して、開始位置を設定する必要あり
		leng += len(d) + 1
	}
	return names
}

func parseTargetVersion(d string) int {

	if len(d) <= 1 {
		return -1
	}

	if d[0] != 'v' {
		return -1
	}

	v := d[1:]
	iv, err := strconv.Atoi(v)
	if err != nil {
		return -1
	}
	return iv
}

func NewModuleGoMod(line string) *Module {

	//package version replace
	s := strings.Split(line, " ")
	str := s[0]

	var m Module
	m.src = str
	m.packageNames = createCandidatePackageNames(str)
	m.current = nil

	//バージョン指定がある場合
	if len(s) >= 2 {
		v := s[1]
		//vを抜いた部分でバージョンを作成
		m.version = NewVersion(v[1:])
		if len(s) >= 3 {
			log.Println("NotSupported replace?:", line)
		}
	}
	return &m
}

func NewModuleManual(str string) *Module {

	var m Module

	m.src = str
	m.packageNames = createCandidatePackageNames(str)
	m.current = nil
	return &m
}

func (m *Module) setCurrent(idx int) error {
	if idx >= len(m.packageNames) {
		return fmt.Errorf("current package name index bound error")
	}
	m.current = m.packageNames[idx]
	return nil
}

// パッケージ名の候補一覧
func (m *Module) Names() []string {
	rtn := make([]string, len(m.packageNames))
	for idx, n := range m.packageNames {
		rtn[idx] = n.name
	}
	return rtn
}

func getModules() ([]*Module, error) {

	lines := callModules()
	if lines == nil {
		return nil, fmt.Errorf("callModules() error: go mod tidy?")
	}

	var mods []*Module
	for _, pkg := range lines {
		m := NewModuleGoMod(pkg)
		if m.version == nil {
			//自パッケージを除去
			continue
		}
		mods = append(mods, m)
	}

	return mods, nil
}
