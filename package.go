package gop

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/xerrors"
)

// TODO バージョン指定のパッケージが存在する場合
// v1を採用するロジックを作成
//
// 全バージョンを確認して、省略形が0か１であることを確認
type Package struct {
	name     string
	path     string
	versions []*Version

	targetVersion *Version
}

func NewPackage(path, n string) *Package {
	var p Package
	p.name = n

	// バージョンを考慮したパスなので注意
	idx := strings.Index(path, "@v")
	p.path = path[:idx]
	p.versions = make([]*Version, 0)
	return &p
}

func (p *Package) AddVersion(v *Version) {
	p.versions = append(p.versions, v)
}

func (p *Package) GetPath() string {

	if p.targetVersion != nil {
		return p.getPath(p.targetVersion)
	}

	p.sortVersion()
	//最後を取得
	v := p.versions[len(p.versions)-1]
	//リリース(PRとBuildのバージョンなし)を確認する
	for idx := len(p.versions) - 1; idx >= 0; idx-- {
		wk := p.versions[idx]
		if wk.IsRelease() {
			v = wk
			break
		}
	}
	return p.getPath(v)
}

func (p *Package) getPath(v *Version) string {
	return fmt.Sprintf("%v@v%v", p.path, v.String())
}

func (p *Package) sortVersion() {
	sort.Slice(p.versions, func(i, j int) bool {
		v1 := p.versions[i]
		v2 := p.versions[j]
		return v1.Lt(v2)
	})
}

func (p *Package) Same(m *Module) bool {

	names := m.Names()
	for idx, name := range names {
		if p.name == name {
			m.setCurrent(idx)
			p.targetVersion = m.version
			return true
		}
	}
	return false
}

func (p *Package) String() string {

	p.sortVersion()
	var b strings.Builder
	b.WriteString(p.name + " Versions:")

	for _, v := range p.versions {
		target := " "
		if p.targetVersion != nil {
			if p.targetVersion.Eq(v) {
				target = "*"
			}
		}
		b.WriteString(fmt.Sprintf("\n  %s %s", target, v.String()))
	}
	return b.String()
}

func (p *Package) Print() {

	p.sortVersion()

	var b strings.Builder
	b.WriteString(p.name + ":")

	for _, v := range p.versions {
		target := " "
		if p.targetVersion != nil {
			if p.targetVersion.Eq(v) {
				target = "*"
			}
		}
		b.WriteString(fmt.Sprintf("\n  %s %s => %s", target, v.String(), p.getPath(v)))
	}
	fmt.Println(b.String())
}

func loadPackages() ([]*Package, error) {

	maps, err := createPackageMaps()
	if err != nil {
		return nil, xerrors.Errorf("createPackageMaps() error: %w", err)
	}

	keys := make([]string, 0, len(maps))
	pkgs := make([]*Package, len(maps))

	for key := range maps {
		keys = append(keys, key)
	}

	sort.Slice(keys, func(i, j int) bool {
		return lessPackage(keys[i], keys[j])
	})

	for idx, key := range keys {
		pkgs[idx] = maps[key]
	}

	return pkgs, nil
}

func createPackageMaps() (map[string]*Package, error) {
	if root == "" {
		return nil, fmt.Errorf("Mod Package Directory Not Found")
	}
	maps := make(map[string]*Package)
	err := setPackages(root, maps)
	if err != nil {
		return nil, fmt.Errorf("setPackage(): %w", err)
	}
	return maps, nil
}

func lessPackage(str1, str2 string) bool {
	return str1 < str2
}

func ignoreDirectory(path string) bool {
	fp := filepath.Join(root, "cache")
	if fp == path {
		return true
	}
	return false
}

func setPackages(path string, pkgs map[string]*Package) error {

	if ignoreDirectory(path) {
		return nil
	}

	list, err := os.ReadDir(path)
	if err != nil {
		return xerrors.Errorf("os.ReadDir() error: %w", err)
	}

	for _, entry := range list {

		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		abs := filepath.Join(path, name)

		n, v, ok := parseVersion(path, name)

		if ok {
			p := pkgs[n]
			if p == nil {
				p = NewPackage(abs, n)
				pkgs[n] = p
			}
			p.AddVersion(v)

		} else {
			err := setPackages(abs, pkgs)
			if err != nil {
				return xerrors.Errorf("setPackages(%s) error: %w", n, err)
			}
		}
	}
	return nil
}

func parseVersion(path, name string) (string, *Version, bool) {

	idx := strings.Index(name, "@v")
	if idx == -1 {
		return "", nil, false
	}

	ver := name[idx+2:]

	v := NewVersion(ver)

	nn := name[:idx]
	p := filepath.Join(path, nn)

	n := strings.ReplaceAll(p, root, "")
	n = strings.ReplaceAll(n, "\\", "/")

	return n[1:], v, true
}

func searchPackage(m *Module) (*Package, error) {

	maps, err := createPackageMaps()
	if err != nil {
		return nil, xerrors.Errorf("createPackageMaps() error: %w", err)
	}

	names := m.Names()
	for _, name := range names {
		p := maps[name]
		if p != nil {
			return p, nil
		}
	}
	return nil, fmt.Errorf("Not Found[%s]", m.src)
}

// 対象となるパッケージを作成する
func removePackages(m *Module, modules []*Module, pkgs []*Package) []*Package {

	//パッケージ指定がない場合
	var mods []*Module = nil
	if modules == nil {
		if m != nil {
			mods = append(mods, m)
		}
	} else {
		if m != nil {
			for _, mod := range modules {
				if m.Eq(mod) {
					mods = append(mods, mod)
					break
				}
			}
			if len(mods) == 0 {
				return nil
			}
		} else {
			mods = modules
		}
	}

	if mods == nil {
		return pkgs
	}

	target := make([]*Package, 0, len(mods))
	for _, mod := range mods {

		var p *Package
		for _, pkg := range pkgs {
			if pkg.Same(mod) {
				p = pkg
				break
			}
		}

		if p == nil {
			if verboseMode {
				log.Println("NotFound:", mod)
			}
		} else {
			target = append(target, p)
		}
	}

	return target
}
