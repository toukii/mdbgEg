package main

import (
	"github.com/everfore/exc"
	// "github.com/shaalx/goutils"
)

var (
	exc_cmd *exc.CMD
)

func init() {
	exc_cmd = exc.NewCMD("rdr").Debug()
	exc_cmd.Execute()
	exc_cmd.Reset("ext -f").Execute()
}

func main() {

}
