package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"net/rpc"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/everfore/exc"
	"github.com/everfore/rpcsv"
	"github.com/toukii/goutils"
	"github.com/toukii/jsnm"
)

func main() {
	defer rpc_client.Close()
	walkRPCRdr()
	upload := &goutils.UploadHandler{}
	mux := http.NewServeMux()
	mux.Handle("/upload", upload)
	mux.Handle("/upload/streamUpload", upload)
	mux.HandleFunc("/callback", callback)
	mux.HandleFunc("/update", update)
	mux.HandleFunc("/TiNews", TiNews)
	mux.HandleFunc("/TiNewsAPI", TiNewsAPI)
	mux.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./MDFs"))))
	http.ListenAndServe(":80", mux)

}

var (
	exc_cmd    *exc.CMD
	rpc_client *rpc.Client
	tpl        *template.Template
	bs         []byte
	// rpc_tcp_server = "localhost:8800"
	// rpc_tcp_server = "tcphub.t0.daoapp.io:61142"
	// rpc_tcp_server = "rpchub.t0.daoapp.io:61142"
	rpc_tcp_server = "rpchub.t1.daoapp.io:61160"
	// rpc_tcp_server = "192.168.1.114:8800"
)

func init() {
	var err error
	exc_cmd = exc.NewCMD("ls").Debug()
	rpc_client = rpcsv.RPCClient(rpc_tcp_server)
	if rpc_client == nil {
		panic("rpc_client is nil!")
	}
	tpl, err = template.ParseFiles("theme.thm")
	if goutils.CheckErr(err) {
		tpl = defaultTheme()
	}
}

func connect() {
	rpc_client = rpcsv.RPCClientWithCodec(rpc_tcp_server)
	go func() {
		time.Sleep(2e9)
		rpc_client.Close()
	}()
}

func defaultTheme() *template.Template {
	dtpl, err := template.New("default").Parse("{{.MDContent}}")
	if goutils.CheckErr(err) {
		panic(err)
	}
	return dtpl
}

func TiNewsAPI(rw http.ResponseWriter, req *http.Request) {
	type TiNews struct {
		Main       string `json:"main"`
		AuthorName string `json:"author_name"`
	}

	news := struct {
		Data []TiNews `json:"data"`
	}{
		Data: []TiNews{
			TiNews{
				Main:       "发布平台运营专员兼软件开发时永宾已确定从银联商务离职。" + time.Now().Format(" 2006/1/2 15:04 "),
				AuthorName: "Ti媒体记者/toukii",
			},
		},
	}
	bs, err := json.Marshal(news)
	if goutils.CheckErr(err) {
		fmt.Fprint(rw, err)
	} else {
		rw.Write(bs)
	}
}

func TiNews(rw http.ResponseWriter, req *http.Request) {
	out := make([]byte, 1)
	i := 0
retry:
	err := rpc_client.Call("RPC.TiNews", &i, &out)
	if goutils.CheckErr(err) {
		if strings.Contains(err.Error(), "connection") && i < 6 {
			connect()
			i++
			goto retry
		}
		fmt.Fprint(rw, err)
		return
	}
	rw.Write(out)
}

func update(rw http.ResponseWriter, req *http.Request) {
	updateTheme()
}

func updateTheme() {
	fmt.Println("update theme")
	exc_cmd.Reset("rm -rf MDFs").Execute()
	tpl1, err := template.ParseFiles("theme.thm")
	if goutils.CheckErr(err) {
		return
	}
	tpl = tpl1
	walkRPCRdr()
}

// Webhooks callback
func callback(rw http.ResponseWriter, req *http.Request) {
	// fmt.Printf("Refer:%s\n", req.Referer())
	// fmt.Printf("req:%#v\n", req)

	usa := req.UserAgent()
	// fmt.Printf("UserAgent:%s\n", usa)
	if !strings.Contains(usa, "GitHub-Hookshot/") && !strings.Contains(usa, "Coding.net Hook") {
		fmt.Println("CSRF Attack!")
		http.Redirect(rw, req, "/", 302)
		return
	}
	/*// coding
	if strings.Contains(usa, "Coding.net Hook") {
		exc_cmd.Reset("git pull origin master:master").Execute()
		rpcsv.UpdataTheme()
		updateTheme()
		return
	}*/
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
		if strings.HasSuffix(fs, ".md") {
			modifiedMD(fs, "./MDFs")
		} else {
			goutils.ReWriteFile(filepath.Join("./MDFs", fs), goutils.ReadFile(fs))
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
		} else {
			goutils.ReWriteFile(filepath.Join("./MDFs", fs), goutils.ReadFile(fs))
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
	i := 0
retry:
	i++
	err = rpc_client.Call("RPC.Markdown", &in, &out)
	// err = rpcsv.Markdown(rpc_client, &in, &out)
	if goutils.CheckErr(err) {
		if i < 6 {
			connect()
			goto retry
		}
		return
	}
	target := fmt.Sprintf("%s.html", filepath.Join(dir_out, dir, fs[0]))
	// goutils.Mkdir(fmt.Sprintf("%s", filepath.Join(dir_out, dir)))
	outfile, erro := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0666)
	if goutils.CheckErr(erro) {
		return
	}
	defer outfile.Close()
	dt := make(map[string]interface{})
	dt["MDContent"] = template.HTML(goutils.ToString(out))
	erre := tpl.Execute(outfile, dt)
	if !goutils.CheckErr(erre) {
		fmt.Println(file_in, " ==> ", target)
	}
}

func copyFile(file_in, dir_out string) {
	goutils.WriteFile(filepath.Join(dir_out, file_in), goutils.ReadFile(file_in))
	fmt.Printf("copy file:%s ==> %s\n", file_in, filepath.Join(dir_out, file_in))
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
	if strings.EqualFold(info.Name(), ".git") || strings.Contains(info.Name(), "MDFs") {
		return filepath.SkipDir
	}
	abspath := filepath.Join(abs, path)
	if info.IsDir() {
		goutils.Mkdir(abspath)
		fmt.Printf("mkdir %s\n", abspath)
		return nil
	}
	copyFile(path, abs)
	/*goutils.ReWriteFile(abspath, goutils.ReadFile(path))
	fmt.Printf("copy file:%s ==> %s\n", path, abspath)*/
	if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
		return nil
	}
	modifiedMD(path, abs)
	return nil
}
