package main

import (
	"github.com/elazarl/goproxy"
	"log"
	"net/http"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true
	proxy.OnRequest(goproxy.DstHostIs("www.bilibili.com")).DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			ctx.Logf("req = %v", r)
			return r, nil
		})
	log.Fatalln(http.ListenAndServe(":8080", proxy))
}