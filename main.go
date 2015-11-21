package main

import (
	"github.com/everfore/exc"
	"github.com/shaalx/goutils"
	"html/template"
	"net/http"
)

var (
	exc_cmd *exc.CMD
)

func init() {
	exc_cmd = exc.NewCMD("rdr").Debug()
}

func main() {
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./MDFs"))))
	http.ListenAndServe(":80", nil)
}

func index(rw http.ResponseWriter, req *http.Request) {
	tpl, err := template.New("README.md.html").ParseFiles("README.md.html")
	if goutils.CheckErr(err) {
		rw.Write(goutils.ReadFile("README.md.html"))
		return
	}
	tpl.Execute(rw, nil)
}
