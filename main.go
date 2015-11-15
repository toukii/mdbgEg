package main

import (
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/shaalx/goutils"
	"io/ioutil"
	"os"
	"text/template"
)

var (
	theme *template.Template
)

const (
	thm_file = "theme.thm"
)

func init() {
	thm_b := readFile(thm_file)
	var err error
	theme, err = template.New("theme.thm").Parse(goutils.ToString(thm_b))
	if goutils.CheckErr(err) {
		panic(err.Error())
	}
}

func main() {
	filename := "README.md"
	renderFile(filename)
}

func renderFile(filename string) bool {
	input := readFile(filename)
	unsafe := blackfriday.MarkdownCommon(input)
	html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)

	data := make(map[string]interface{})
	data["MDContent"] = goutils.ToString(html)
	data["Title"] = filename
	of, err := os.OpenFile(filename+".html", os.O_CREATE|os.O_RDWR, 0666)
	defer of.Close()
	if goutils.CheckErr(err) {
		return false
	}
	err = theme.Execute(of, data)
	if goutils.CheckErr(err) {
		return false
	}
	return true
}

func readFile(filename string) []byte {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0644)
	if goutils.CheckErr(err) {
		return nil
	}
	defer file.Close()
	b, err := ioutil.ReadAll(file)
	if goutils.CheckErr(err) {
		return nil
	}
	return b
}
