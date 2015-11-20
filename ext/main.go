package main

import (
	"flag"
	"fmt"
	"github.com/everfore/exc"
	"github.com/shaalx/goutils"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	Spor  = string(os.PathSeparator)
	from  = "./" // extration files in from
	home  = ""   // home path == pwd + Spor + base
	base  = ""   // base path
	ext   = ".md"
	force = false

	exc_cmd *exc.CMD
)

func init() {
	_home, _ := os.Getwd()
	// flag.StringVar(&from, "from", "./", "-from MDFs")
	flag.StringVar(&base, "base", "MDFs", "-base mdfs")
	flag.StringVar(&ext, "ext", ".md", "-ext .go")
	flag.BoolVar(&force, "f", false, "-f [true] force")
	flag.Parse()
	home = _home + Spor + base
	exc_cmd = exc.NewCMD("ls").Debug()
}

func main() {
	fmt.Printf("Extraction base path:%s.\n", base)
	_, err := os.Lstat(home)
	if err == nil {
		if force {
			err = os.RemoveAll(base)
			if goutils.CheckErr(err) {
				return
			}
		} else {
			fmt.Printf("%s has existed!", home)
			return
		}
	}
	Extraction(from)
	filepath.Walk(from, RenderFunc)
}

func ExtractFunc(path string, info os.FileInfo, err error) error {
	if strings.EqualFold(".git", info.Name()) {
		return filepath.SkipDir
	}
	if strings.HasPrefix(path, base) {
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
		fmt.Printf("%s: %d bytes.\n", path, n)
	}
	return nil
}

func RenderFunc(path string, info os.FileInfo, err error) error {
	if strings.EqualFold(".git", info.Name()) {
		return filepath.SkipDir
	}
	if strings.HasPrefix(path, base) {
		return filepath.SkipDir
	}
	if !info.IsDir() && strings.EqualFold(ext, filepath.Ext(path)) {
		dir := filepath.Dir(path)
		if dir == "." {
			dir = ""
		}
		// fromd := strings.Replace(home+Spor+path, "/", "\\\\", -1)
		// tod := strings.Replace(home+Spor+dir, "/", "\\\\", -1)
		// cmd := fmt.Sprintf("md -r -f %s -d %s", fromd, tod)
		cmd := fmt.Sprintf("md -r -f %s -d %s", home+Spor+path, home+Spor+dir)
		fmt.Println(cmd)
		exc_cmd.Reset(cmd).Wd().Execute()
	}
	return nil
}

func Extraction(from string) {
	filepath.Walk(from, ExtractFunc)
}
