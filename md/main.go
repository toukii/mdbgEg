package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	// "github.com/microcosm-cc/bluemonday"
	// "github.com/russross/blackfriday"
	"github.com/shaalx/goutils"
	md "github.com/shurcooL/github_flavored_markdown"
	"os"
	"strings"
	"text/template"
)

var (
	filename  = ""
	targetDir = ""
	redo      = false
	theme     *template.Template
	Spor      = ""
)

const (
	thm_file = "theme.thm"
	thm_s    = `<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<title>{{.Title}}</title>
	<link href="/favicon.png" rel="icon">

	<link href="http://7xku3c.com1.z0.glb.clouddn.com/md_style.css" rel="stylesheet">

	<link href="http://cdn.bootcss.com/bootstrap/3.3.4/css/bootstrap.min.css" rel="stylesheet">
    <link href="http://cdn.bootcss.com/font-awesome/4.2.0/css/font-awesome.min.css" rel="stylesheet">
    <link href="http://static.bootcss.com/www/assets/css/site.min.css?v5" rel="stylesheet">
    <link crossorigin="anonymous" href="https://assets-cdn.github.com/assets/github2-53964e9b93636aa437196c028e3b15febd3c6d5a52d4e8368a9c2894932d294e.css" integrity="sha256-U5ZOm5NjaqQ3GWwCjjsV/r08bVpS1Og2ipwolJMtKU4=" media="all" rel="stylesheet" />
</head>
	<body>
		<div class="container">
			<nav class="navbar navbar-default" role="navigation" id="navbar">
				<div class="collapse navbar-collapse navbar-ex1-collapse">
					<ul class="nav navbar-nav" id="menu">
						<li><a href="/">Home</a></li>
						<li><a href="/Go">Go</a></li>
						<li><a href="/Item">Item</a></li>
					</ul>
				</div>
			</nav>
		</div>

		<div class="container">
            <div class="col-md-8">			
				<div class="content">
					{{.MDContent}}
				</div>
		</div>

			<div class="col-md-4 sidebar">
			  <div class="panel panel-default">
				<div class="panel-body">
				  <div align="left">
				  	<h4><small>学习链接</small></h4>
				  </div>
				  <hr>					
					<strong><a href="https://gowalker.org/" title="gowalker" rel="nofollow">gowalker</a></strong> 
					<strong><a href="https://godoc.org/" title="godoc" rel="nofollow">godoc</a></strong> 
					<strong><a href="https://gopm.io/" title="gopm" rel="nofollow">gopm</a></strong> 
					<strong><a href="http://stdlib-shaalx.myalauda.cn/" title="stdlib" rel="nofollow">gostd</a></strong>
				</div>
			  </div>

<div class="panel panel-default">
	<div class="panel-heading">
	  <h3 class="panel-title">状态</h3>
	</div>
	<table width="100%" class="status">
	  <thead>
		<tr>
		  <th>&nbsp;</th>
		  <th></th>
		</tr>
	  </thead>
	  <tbody>
		<tr>
		  <td class="status-label">Go</td>
		  <td class="value">4</td>
		</tr>
		<tr>
		  <td class="status-label">Java</td>
		  <td class="value">0</td>
		</tr>
		<tr>
		  <td class="status-label">Others</td>
		  <td class="value">0</td>
		</tr>
	  </tbody>
	</table>
  </div>

		</div>

		<div class="col-md-12">
		<footer class="footer">
			<div class="row footer-bottom">
				<ul class="list-inline text-center">
					<div class="copy-right" style="color:#4d5152">
						<h6><small> ©2015 shaalx </small></h6>
					</div>
				</ul>
			</div>
		</footer>
		</div>
	</body>
</html>`
)

func init() {
	Spor = string(os.PathSeparator)
	// flag
	flag.BoolVar(&redo, "r", false, "-r [true]")
	flag.StringVar(&filename, "f", "README.md", "-f readme.md")
	flag.StringVar(&targetDir, "d", "./", "-d ./static")

	// theme
	// thm_b := readFile(thm_file)
	// if nil == thm_b {
	thm_b := goutils.ToByte(thm_s)
	// }
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
	_, err := os.Lstat(filename + ".html")
	if !redo && nil == err {
		return false
	}
	if nil == err {
		err = os.Remove(filename + ".html")
		goutils.CheckErr(err)
	}
	input := readFile(filename)
	if nil == input {
		return false
	}
	// unsafe := blackfriday.MarkdownCommon(input)
	// html := bluemonday.UGCPolicy().SanitizeBytes(unsafe)
	html := md.Markdown(input)

	data := make(map[string]interface{})
	data["MDContent"] = goutils.ToString(html)

	data["Title"] = getName(filename)
	of, err := os.OpenFile( /*targetDir+string(os.PathSeparator)+*/ filename+".html", os.O_CREATE|os.O_RDWR, 0666)
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

func getName(filename string) string {
	sps := strings.Split(filename, Spor)
	return sps[len(sps)-1]
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
