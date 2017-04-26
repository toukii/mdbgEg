package main

import (
	"fmt"
	"github.com/everfore/exc"
	"net/http"
	"strings"

	"github.com/toukii/jsnm"

	"html/template"
	"os"
	"path/filepath"

	"github.com/everfore/rpcsv"
	"github.com/toukii/goutils"
	// "net"
	"net/rpc"
)

func main() {
	defer rpc_client.Close()
	// defer lis.Close()
	walkRPCRdr()
	copyFile("index.html", "./MDFs")
	http.HandleFunc("/callback", callback)
	http.HandleFunc("/update", update)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./MDFs"))))
	http.ListenAndServe(":80", nil)
}

var (
	exc_cmd    *exc.CMD
	rpc_client *rpc.Client
	// lis        net.Listener
	tpl *template.Template
)

func init() {
	var err error
	// if lis, err = rpcsv.RPCServe("88"); err != nil {
	// 	return
	// }
	exc_cmd = exc.NewCMD("ls").Debug()
	// rpc_client = rpcsv.RPCClient("127.0.0.1:88")
	rpc_client = rpcsv.RPCClient("tcphub.t0.daoapp.io:61142")
	if rpc_client == nil {
		// lis.Close()
		panic("rpc_client is nil!")
	}
	tpl, err = template.ParseFiles("theme.thm")
	if goutils.CheckErr(err) {
		tpl = defaultTheme()
	}
}

func defaultTheme() *template.Template {
	dtpl, err := template.New("default").Parse("{{.MDContent}}")
	if goutils.CheckErr(err) {
		panic(err)
	}
	return dtpl
}

func update(rw http.ResponseWriter, req *http.Request) {
	updateTheme()
}

func updateTheme() {
	fmt.Println("update")
	tpl1, err := template.ParseFiles("theme.thm")
	if goutils.CheckErr(err) {
		return
	}
	tpl = tpl1
	walkRPCRdr()
}

// Webhooks callback
func callback(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Refer:%s\n", req.Referer())
	fmt.Printf("req:%#v\n", req)

	usa := req.UserAgent()
	// fmt.Printf("UserAgent:%s\n", usa)
	/*if !strings.Contains(usa, "GitHub-Hookshot/") && !strings.Contains(usa, "Coding.net Hook") {*/
	/*	fmt.Println("CSRF Attack!")*/
	/*	http.Redirect(rw, req, "/", 302)*/
	/*	return*/
	/*}*/
	// coding
	if strings.Contains(usa, "Coding.net Hook") {
		exc_cmd.Reset("git pull origin master:master").Execute()
		rpcsv.UpdataTheme()
		updateTheme()
		return
	}
	// coding
	hj := jsnm.ReaderFmt(req.Body)
	ma := hj.ArrGet("commits", "0", "modified").Arr()
	pull := false
	if len(ma) > 0 {
		exc_cmd.Reset("git pull origin master:master").Execute()
		pull = true
	}
	for i, it := range ma {
		fs := it.RawData().String()
		fmt.Printf("modified-%d:%v\n", i, fs)
		if strings.EqualFold(fs, "theme.thm") {
			rpcsv.UpdataTheme()
			updateTheme()
			return
		}
		if strings.HasSuffix(fs, ".html") {
			goutils.WriteFile(filepath.Join("./MDFs", fs), goutils.ReadFile(fs))
		}
		if strings.HasSuffix(fs, ".md") {
			modifiedMD(fs, "./MDFs")
		}
	}
	aa := hj.ArrGet("commits", "0", "added").Arr()
	if aa != nil && !pull {
		exc_cmd.Reset("git pull origin master:master").Execute()
	}
	for i, it := range aa {
		fs := it.RawData().String()
		fmt.Printf("added-%d:%v\n", i, fs)
		if strings.HasSuffix(fs, ".md") {
			modifiedMD(fs, "./MDFs")
		}
	}
	ra := hj.ArrGet("commits", "0", "removed").Arr()
	if ra != nil && !pull {
		exc_cmd.Reset("git pull origin master:master").Execute()
	}
	for i, it := range ra {
		fs := it.RawData().String()
		fmt.Printf("removed-%d:%v\n", i, fs)
		if strings.HasSuffix(fs, ".md") {
			removeMD(fs, "./MDFs")
		}
	}
}

func removeMD(file_in, dir_out string) {
	fs := strings.Split(file_in, ".")
	goutils.DeleteFile(fmt.Sprintf("%s.html", filepath.Join(dir_out, fs[0])))
}

// in: Linux/index.md
// out: ./MDFs
func modifiedMD(file_in, dir_out string) {
	finfo, err := os.Stat(file_in)
	if goutils.CheckErr(err) {
		return
	}
	filename := finfo.Name()
	dir := filepath.Dir(file_in)
	fs := strings.Split(filename, ".")
	in := goutils.ReadFile(file_in)
	out := make([]byte, 1)
	err = rpcsv.Markdown(rpc_client, &in, &out)
	if goutils.CheckErr(err) {
		return
	}
	// fmt.Println(goutils.ToString(out))
	target := fmt.Sprintf("%s.html", filepath.Join(dir_out, dir, fs[0]))

	goutils.Mkdir(fmt.Sprintf("%s", filepath.Join(dir_out, dir)))
	outfile, _ := os.OpenFile(fmt.Sprintf("%s.html", filepath.Join(dir_out, dir, fs[0])), os.O_CREATE|os.O_WRONLY, 0666)
	defer outfile.Close()
	dt := make(map[string]interface{})
	dt["MDContent"] = template.HTML(goutils.ToString(out))
	tpl.Execute(outfile, dt)
	fmt.Println(file_in, " ==> ", target)
}

func copyFile(file_in, dir_out string) {
	goutils.WriteFile(filepath.Join(dir_out, file_in), goutils.ReadFile(file_in))
}

// base: ./
// target: ./MDFs
func walkRPCRdr() {
	filepath.Walk("./", walkCond)
}

var (
	abs, _ = filepath.Abs("./MDFs")
)

func walkCond(path string, info os.FileInfo, err error) error {
	if strings.EqualFold(info.Name(), ".git") || strings.EqualFold(info.Name(), "./MDFs") {
		return filepath.SkipDir
	}
	if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") || !strings.HasSuffix(info.Name(), ".html") {
		return nil
	}
	modifiedMD(path, abs)
	return nil
}
