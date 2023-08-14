package gop

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

var (
	root        string = ""
	verboseMode bool   = false
)

func init() {

	p := getGOPATH()
	if p == "" {
		log.Println("GOPATH(go env GOPATH) not exists")
		return
	}
	modp := filepath.Join(p, "pkg", "mod")
	if _, err := os.Stat(modp); err != nil {
		log.Println(err)
		modp = ""
	}
	root = modp
}

func Print(verbose, project, all, list bool, input string) error {

	verboseMode = verbose
	if verboseMode {
		printGoVersion()
	}

	if root == "" {
		return fmt.Errorf(`"go" does not exist`)
	}

	var err error
	//プロジェクトのモジュール一覧
	var modules []*Module
	//指定のモジュール
	var m *Module
	//実パッケージ位置
	var target []*Package

	//実パッケージを取得
	pkgs, err := loadPackages()
	if err != nil {
		return fmt.Errorf("loadPackages() error: %w", err)
	}

	if project {
		//go.modからモジュールを取得
		modules, err = getModules()
		if err != nil {
			return fmt.Errorf("getModules() error: %w", err)
		}
	}

	//指定がない場合
	if input != "" {
		//入力からモジュールを取得
		m = NewModuleManual(input)
		if verboseMode {
			fmt.Println(m)
		}
	}

	target = removePackages(m, modules, pkgs)
	//  対象パッケージをすべて表示
	for _, pkg := range target {
		if list {
			pkg.Print()
		} else {
			//TODO 指定方法
			fmt.Println(pkg.GetPath())
		}
	}
	return nil
}

func printGoVersion() {
	out, err := cmd("go", "version")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(out)
}

func getGOPATH() string {
	out, err := cmd("go", "env", "GOPATH")
	if err != nil {
		log.Printf("%+v", err)
		return ""
	}
	return out
}

func cmd(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.Output()
	if err != nil {
		return "", xerrors.Errorf("command error: %w", err)
	}
	return strings.Trim(string(out), "\r\n"), nil
}

func parseInt(buf string) int {
	v, err := strconv.Atoi(buf)
	if err != nil {
		return -1
	}
	return v
}

func findGoMod() string {
	out, err := cmd("go", "env", "GOMOD")
	if err != nil {
		log.Printf("%+v", err)
	}
	return out
}

func callModules() []string {
	out, err := cmd("go", "list", "-m", "all")
	if err != nil {
		log.Printf("%+v", err)
	}

	rtn := strings.Split(out, "\n")

	for idx, line := range rtn {
		rtn[idx] = strings.Trim(line, "\r\n")
	}

	return rtn
}
