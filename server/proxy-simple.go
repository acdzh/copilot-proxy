package main

import (
	// "crypto/tls"

	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	// 如果你需要代理, 使用 transport 相关的这几行
	// // create transPort.
	// clashUrl, _ := url.Parse("http://127.0.0.1:7890")
	// transport := &http.Transport{
	// 	Proxy:           http.ProxyURL(clashUrl),
	// 	TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 跳过 HTTPS 证书验证
	// }

	// create reverse proxy
	target, _ := url.Parse("https://copilot-proxy.githubusercontent.com") // 目标 API 的地址
	proxy := httputil.NewSingleHostReverseProxy(target)
	// proxy.Transport = transport

	// create http server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy.ServeHTTP(w, r) // 将请求转发到目标 API
	})
	log.Println("server start at http://127.0.0.1:9394")
	http.ListenAndServe(":9394", nil)
}
