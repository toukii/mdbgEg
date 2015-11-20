package main

import (
	"flag"
	"fmt"
	"github.com/shaalx/goutils"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	Spor  = string(os.PathSeparator)
	home  = "" // home path == pwd + Spor + base
	base  = "" // base path
	ext   = ".md"
	force = false
)

func init() {
	_home, _ := os.Getwd()
	flag.StringVar(&base, "b", "MDFs", "-b mdfs")
	flag.StringVar(&ext, "e", ".md", "-e .go")
	flag.BoolVar(&force, "f", false, "-f [true]")
	flag.Parse()
	home = _home + Spor + base
}

func main() {
	fmt.Println(base)
	_, err := os.Lstat(home)
	if !goutils.CheckErr(err) {
		if force {
			os.Remove(base)
		} else {
			fmt.Printf("%s has existed!", home)
		}
	}
	Extraction("./")
}

func WalkFunc(path string, info os.FileInfo, err error) error {
	if strings.EqualFold(".git", info.Name()) {
		return filepath.SkipDir
	}
	if strings.Contains(path, base) {
		return filepath.SkipDir
	}
	if info.IsDir() {
		err = os.Mkdir(home+Spor+path, 0644)
		if goutils.CheckErr(err) {
			return nil
		}
	}

	if strings.EqualFold(ext, filepath.Ext(path)) {
		orf, err := os.OpenFile(path, os.O_RDONLY, 0644)
		defer orf.Close()
		if goutils.CheckErr(err) {
			return nil
		}
		owf, err := os.OpenFile(base+string(os.PathSeparator)+path, os.O_CREATE|os.O_WRONLY, 0622)
		defer owf.Close()
		if goutils.CheckErr(err) {
			return nil
		}
		n, err := io.Copy(owf, orf)
		fmt.Printf("n:%d,err:%v\n", n, err)
	}
	return nil
}

func Extraction(base string) {
	filepath.Walk("./", WalkFunc)
}
