// SimpleFileClient project main.go
package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
)

var Url = flag.String("u", "http://120.27.196.178:54321/", "-u=<url> default=http://120.27.196.178:54321/")
var Method = flag.String("m", "upload", "-m=<upload | cmd | download> default=upload")
var Dst = flag.String("d", "", "-d=<dst file> 如果 -m=upload|download 此参数必须填写")
var Src = flag.String("s", "", "-s=<src file> 如果 -m=upload|download 此参数必须填写")
var CmdN = flag.String("c", "", "-c=<cmd> 如果 -m=cmd 此参数必须填写")
var Arg = flag.String("a", "[]", "-a=<arg> 如果 -a=[] 此参数必须填写")

func Check(resp *http.Response) {
	Code := resp.Header.Get("Code")
	Message := resp.Header.Get("Message")
	if Code != "0" && Code != "" {
		panic(errors.New("code:" + Code + ";message:" + Message))
	}
}

func getReaderio(s string) io.ReadCloser {
	if s[0:4] == "http" {
		r, err := http.Get(s)
		if err != nil {
			panic(errors.New(err.Error()))
		}
		return r.Body
	}

	file, err := os.OpenFile(s, os.O_RDONLY, 0666)
	if err != nil {
		panic(errors.New(err.Error()))
	}
	return file
}

func Upload() {
	sio := getReaderio(*Src)
	defer sio.Close()
	r := bufio.NewReader(sio)
	resp, err := http.Post(*Url+*Method+"?dst="+*Dst, "text/plain", r)
	if err != nil {
		panic(errors.New(err.Error()))
	}
	defer resp.Body.Close()
	Check(resp)
	w := bufio.NewWriter(os.Stdout)
	io.Copy(w, resp.Body)
	w.Flush()
}

func Download() {
	resp, err := http.Post(*Url+*Method+"?src="+*Src, "text/plain", nil)
	if err != nil {
		panic(errors.New(err.Error()))
	}
	defer resp.Body.Close()
	Check(resp)
	file, err := os.OpenFile(*Dst, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(errors.New(err.Error()))
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	io.Copy(w, resp.Body)
	w.Flush()
}

func Cmd() {
	u := *Url + *Method + "?cmd=" + *CmdN + "&arg=" + url.PathEscape(*Arg)
	fmt.Println("CMD POST", u)
	resp, err := http.Post(u, "text/plain", nil)
	if err != nil {
		panic(errors.New(err.Error()))
	}
	defer resp.Body.Close()
	Check(resp)
	w := bufio.NewWriter(os.Stdout)
	io.Copy(w, resp.Body)
	w.Write([]byte("\n"))
	w.Flush()
}

func main() {
	flag.Parse()
	if len(*Url) == 0 {
		panic(errors.New("url 不能为空"))
	}
	if (*Url)[len(*Url)-1] != '/' {
		*Url = *Url + "/"
	}
	if *Method == "upload" {
		Upload()
	} else if *Method == "cmd" {
		Cmd()
	} else {
		panic(errors.New(*Method + " 方法不支持"))
	}
}
