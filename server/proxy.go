package main

import (
	"bytes"
	// "crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

func colorRed(s string) string {
	return "\033[31m" + s + "\033[0m"
}

func colorGreen(s string) string {
	return "\033[32m" + s + "\033[0m"
}

func colorYellow(s string) string {
	return "\033[33m" + s + "\033[0m"
}

func colorBlue(s string) string {
	return "\033[34m" + s + "\033[0m"
}

func colorMagenta(s string) string {
	return "\033[35m" + s + "\033[0m"
}

func colorCyan(s string) string {
	return "\033[36m" + s + "\033[0m"
}

func colorWhite(s string) string {
	return "\033[37m" + s + "\033[0m"
}

func colorGray(s string) string {
	return "\033[90m" + s + "\033[0m"
}

// 定义一个自定义的 ResponseWriter 对象
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

// 实现 Write 方法，拦截所有写入到 ResponseWriter 的数据
func (w *loggingResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// 实现 WriteHeader 方法，拦截设置的状态码
func (w *loggingResponseWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func main() {
	// 如果你需要代理, 使用 transport 相关的这几行
	// create transPort.
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
		start := time.Now()
		ip := strings.Split(r.RemoteAddr, ":")[0]

		reqBodyStr := func() string {
			if r.ContentLength > 0 {
				reqBodyBytes, err := ioutil.ReadAll(r.Body)
				if err != nil {
					log.Fatal(err)
					return "[no body]"
				}
				r.Body = ioutil.NopCloser(bytes.NewBuffer(reqBodyBytes))
				return string(reqBodyBytes)
			} else {
				return "[no body]"
			}
		}()

		// 创建自定义的 ResponseWriter 对象
		lrw := &loggingResponseWriter{w, 0, new(bytes.Buffer)}
		proxy.ServeHTTP(lrw, r) // 将请求转发到目标 API

		resBodyStr := func() string {
			if lrw.Header().Get("Content-Length") == "0" {
				return "[no body]"
			} else {
				return lrw.body.String()
			}
		}()

		log.Printf("\n")
		log.Printf("%s %s %s to %s from %s, Completed in %v", colorYellow("--->"), func() string {
			if lrw.statusCode == 200 {
				return colorGreen(fmt.Sprintf("[%d]", lrw.statusCode))
			} else {
				return colorRed(fmt.Sprintf("[%d]", lrw.statusCode))
			}
		}(), colorBlue(fmt.Sprintf("[%s]", r.Method)), colorYellow(r.URL.Path), colorBlue(ip), colorBlue(time.Since(start).String()))
		log.Printf("     %s: %s", colorBlue("[header]"), colorGray(fmt.Sprintf("%v", r.Header)))
		if reqBodyStr != "[no body]" {
			log.Printf("     %s: %s", colorBlue("[body]"), colorGray(reqBodyStr))
		}
		log.Printf("%s %s: %s", colorYellow("<---"), colorBlue("[body]"), colorGray(strings.ReplaceAll(resBodyStr, "\n\n", "\n")))
	})
	log.Println("server start at http://127.0.0.1:9394")
	http.ListenAndServe(":9394", nil)

	// if you want use gin
	// router := gin.Default()
	// router.Any("/*path", func(c *gin.Context) {
	// 	proxy.ServeHTTP(c.Writer, c.Request) // 将请求转发到目标 API
	// })
	// router.Run(":9394")
}
