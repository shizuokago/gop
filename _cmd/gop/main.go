package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/shizuokago/gop"
	"golang.org/x/xerrors"
)

var (
	verboseMode bool
	allMode     bool
	listMode    bool
)

func init() {
	flag.BoolVar(&verboseMode, "v", false, "Verbose mode")
	flag.BoolVar(&allMode, "all", false, "All mode")
	flag.BoolVar(&listMode, "list", false, "Package List mode")
}

/**
 * Go modules Path
 *
 * gocd
 * mod > ls
 * all package name
 *
 * gocd golang.org/x/xerrors
 * -> $GOPATH/pkg/mod/
 *
 * @マークがあればパッケージ
 */
func main() {

	flag.Parse()
	args := flag.Args()
	err := run(args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "run() error:\n%+v\n", err)
		os.Exit(1)
	}
	return
}

func run(args []string) error {

	leng := len(args)
	pkg := ""
	if leng != 0 {
		pkg = args[0]
	}

	err := gop.Print(verboseMode, allMode, listMode, pkg)
	if err != nil {
		return xerrors.Errorf("gocd.Print() error: %w", err)
	}
	return nil
}
