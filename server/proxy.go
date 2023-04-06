package main

import (
	"bytes"
	"io"
	"math/rand"
	"os"

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

func getBodyFromHttpRequest(r *http.Request) string {
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
}

type ServerLogParams struct {
	statusCode int
	method     string
	path       string
	ip         string
	reqHeader  map[string][]string
	startTime  time.Time
	reqBodyStr string
	resBodyStr string
}

func serverLog(p ServerLogParams) {
	log.Printf("\n")
	log.Printf("%s %s %s to %s from %s, Completed in %v", colorYellow("--->"), func() string {
		if p.statusCode == 200 {
			return colorGreen(fmt.Sprintf("[%d]", p.statusCode))
		} else {
			return colorRed(fmt.Sprintf("[%d]", p.statusCode))
		}
	}(), colorBlue(fmt.Sprintf("[%s]", p.method)), colorYellow(p.path), colorBlue(p.ip), colorBlue(time.Since(p.startTime).String()))
	log.Printf("     %s: %s", colorBlue("[header]"), colorGray(fmt.Sprintf("%v", p.reqHeader)))
	if p.reqBodyStr != "[no body]" {
		log.Printf("     %s: %s", colorBlue("[body]"), colorGray(p.reqBodyStr))
	}
	log.Printf("%s %s: %s", colorYellow("<---"), colorBlue("[body]"), colorGray(strings.ReplaceAll(p.resBodyStr, "\n\n", "\n")))
}

func setUpLogger() {
	// 创建一个文件，用于写入日志信息
	logFile, err := os.OpenFile("copilot-proxy.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file: ", err)
	}
	// defer logFile.Close()
	// 创建一个 MultiWriter，用于同时将日志信息输出到标准错误流和文件中
	mw := io.MultiWriter(os.Stdout, logFile)
	// 设置 log 包的输出流为 MultiWriter
	log.SetOutput(mw)
}

func main() {
	setUpLogger()

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
		reqBodyStr := getBodyFromHttpRequest(r)

		// 不样扫描
		if strings.Contains(r.Header.Get("User-Agent"), "SecurityScan-PoCScan") ||
			r.Header.Get("X-Scan") != "" ||
			r.Header.Get("X-Scan-Active") != "" ||
			r.Header.Get("X-Scan-Token") != "" {
			statusCode, resBodyStr := func() (int, string) {
				random := rand.Intn(100)
				if random <= 45 {
					return http.StatusTeapot, "[418] Stop scanning, I'm a teapot 🫖."
				} else if random <= 80 {
					return http.StatusOK, "[200] Stop scanning, I'm not a teapot 🫖."
				} else {
					return http.StatusOK, "[???] Stop scanning, I'm an eggplant 🍆."
				}
			}()

			w.WriteHeader(statusCode)
			w.Write([]byte(resBodyStr))

			serverLog(ServerLogParams{
				statusCode: statusCode,
				method:     r.Method,
				path:       r.URL.Path,
				ip:         ip,
				reqHeader:  r.Header,
				startTime:  start,
				reqBodyStr: reqBodyStr,
				resBodyStr: resBodyStr,
			})
			return
		}

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

		serverLog(ServerLogParams{
			statusCode: lrw.statusCode,
			method:     r.Method,
			path:       r.URL.Path,
			ip:         ip,
			reqHeader:  r.Header,
			startTime:  start,
			reqBodyStr: reqBodyStr,
			resBodyStr: resBodyStr,
		})
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
