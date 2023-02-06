package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type MyServer struct {
	root string

	fileServer   http.Handler
	reverseSever *httputil.ReverseProxy
}

func (f MyServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upath := r.URL.Path
	if !strings.HasPrefix(upath, "/") {
		upath = "/" + upath
	}
	_, err := os.Stat(filepath.Join(f.root, upath))
	if os.IsNotExist(err) {
		f.reverseSever.ServeHTTP(w, r)
	} else {
		//r.URL.Path = upath
		f.fileServer.ServeHTTP(w, r)
	}
}

func NewMyServer(root string, proxyURL string) MyServer {
	var s MyServer
	s.root = root

	s.fileServer = http.FileServer(http.Dir(root))
	backendURL, err := url.Parse(proxyURL)
	if err != nil {
		fmt.Printf("backend proxy server invalid: %v\n", err)
	}

	s.reverseSever = httputil.NewSingleHostReverseProxy(backendURL)

	return s
}
