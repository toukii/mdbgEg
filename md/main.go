package main

import (
	"flag"
	"fmt"
	"github.com/microcosm-cc/bluemonday"
	"github.com/russross/blackfriday"
	"github.com/shaalx/goutils"
	"io/ioutil"
	"os"
	"text/template"
)

var (
	filename = ""
	redo     = false
	theme    *template.Template
)

const (
	thm_file = "theme.thm"
	thm_s    = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>{{.Title}}</title>
	<!-- <link href="http://cdn.bootcss.com/bootstrap/3.3.4/css/bootstrap.min.css" rel="stylesheet">
    <link href="http://cdn.bootcss.com/font-awesome/4.2.0/css/font-awesome.min.css" rel="stylesheet">
    <link href="http://static.bootcss.com/www/assets/css/site.min.css?v5" rel="stylesheet"> -->
    <link crossorigin="anonymous" href="https://assets-cdn.github.com/assets/github-6670887f84dea33391b25bf5af0455816ab82a9bec8f4f5e4d8160d53b08c0f3.css" integrity="sha256-ZnCIf4TeozORslv1rwRVgWq4Kpvsj09eTYFg1TsIwPM=" media="all" rel="stylesheet" />
    <link crossorigin="anonymous" href="https://assets-cdn.github.com/assets/github2-53964e9b93636aa437196c028e3b15febd3c6d5a52d4e8368a9c2894932d294e.css" integrity="sha256-U5ZOm5NjaqQ3GWwCjjsV/r08bVpS1Og2ipwolJMtKU4=" media="all" rel="stylesheet" />
    {{.MoreCss}}
</head>
	<body>
		<div class="container">
			{{.MDContent}}
		</div>
	</body>
</html>`
)

func init() {
	// flag
	flag.BoolVar(&redo, "r", false, "-r [true]")
	flag.StringVar(&filename, "f", "README.md", "-f readme.md")

	// theme
	// thm_b := readFile(thm_file)
	thm_b := goutils.ToByte(thm_s)
	var err error
	theme, err = template.New("theme.thm").Parse(goutils.ToString(thm_b))
	if goutils.CheckErr(err) {
		panic(err.Error())
	}
}

func main() {
	flag.Parse()
	if renderFile(filename, redo) {
		fmt.Printf("Successfully parsed %s ==> %s.html\n", filename, filename)
	} else {
		fmt.Printf("Failed to parse %s ==> %s.html\n", filename, filename)
	}
}

func renderFile(filename string, redo bool) bool {
	_, err := os.Stat(filename + ".html")
	if !redo && nil == err {
		return false
	}
	input := readFile(filename)
	if nil == input {
		return false
	}
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
