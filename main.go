package main

import (
	"fmt"
	"github.com/everfore/exc"
	"net/http"
	"strings"

	"github.com/shaalx/jsnm"

	"os"
	"path/filepath"

	"github.com/everfore/rpcsv"
	"github.com/everfore/rpcsv/clt"
	"github.com/shaalx/goutils"
	"net"
	"net/rpc"
)

func main() {
	defer rpc_client.Close()
	defer lis.Close()
	walkRPCRdr()
	http.HandleFunc("/callback", callback)
	http.Handle("/", http.StripPrefix("/", http.FileServer(http.Dir("./MDFs"))))
	http.ListenAndServe(":80", nil)
}

var (
	exc_cmd    *exc.CMD
	rpc_client *rpc.Client
	lis        net.Listener
)

func init() {
	var err error
	if lis, err = rpcsv.RPCServe("88"); err != nil {
		return
	}
	exc_cmd = exc.NewCMD("ls").Debug()
	rpc_client = clt.RPCClient("127.0.0.1:88")
	if rpc_client == nil {
		lis.Close()
		panic("rpc_client is nil!")
	}
}

// Webhooks callback
func callback(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Refer:%s\n", req.Referer())
	fmt.Printf("req:%#v\n", req)

	usa := req.UserAgent()
	fmt.Printf("UserAgent:%s\n", usa)
	if !strings.Contains(usa, "GitHub-Hookshot/") {
		fmt.Println("CSRF Attack!")
		http.Redirect(rw, req, "/", 302)
		return
	}
	hj := jsnm.ReaderFmt(req.Body)
	ma := hj.Get("commits").ArrLoc(0).Get("modified").Arr()
	pull := false
	if len(ma) > 0 {
		exc_cmd.Reset("git pull origin master:master").Execute()
		pull = true
	}
	for i, it := range ma {
		fs := it.RawData().String()
		fmt.Printf("modified-%d:%v\n", i, fs)
		if strings.HasSuffix(fs, ".md") {
			modifiedMD(fs, "./MDFs")
		}
	}
	aa := hj.Get("commits").ArrLoc(0).Get("added").Arr()
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
	ra := hj.Get("commits").ArrLoc(0).Get("removed").Arr()
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
	err = clt.Markdown(rpc_client, &in, &out)
	if goutils.CheckErr(err) {
		return
	}
	target := fmt.Sprintf("%s.html", filepath.Join(dir_out, dir, fs[0]))
	fmt.Println(file_in, " ==> ", target)
	goutils.WriteFile(fmt.Sprintf("%s.html", filepath.Join(dir_out, dir, fs[0])), out)
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
	if info.IsDir() || !strings.HasSuffix(info.Name(), ".md") {
		return nil
	}
	modifiedMD(path, abs)
	return nil
}
