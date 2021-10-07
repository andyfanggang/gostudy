package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/felixge/httpsnoop"
)

///定义一个结构体实现http.handle接口
type getHealthzHandler struct{}

func (h getHealthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "<h1>200</h1>")
}

func ipAddrFromRemoteAddr(s string) string {
	idx := strings.LastIndex(s, ":")
	if idx == -1 {
		return s
	}
	return s[:idx]
}

// requestGetRemoteAddress 返回发起请求的客户端 ip 地址，这是出于存在 http 代理的考量
func requestGetRemoteAddress(r *http.Request) string {
	hdr := r.Header
	hdrRealIP := hdr.Get("X-Real-Ip")
	hdrForwardedFor := hdr.Get("X-Forwarded-For")
	if hdrRealIP == "" && hdrForwardedFor == "" {
		return ipAddrFromRemoteAddr(r.RemoteAddr)
	}
	if hdrForwardedFor != "" {
		// X-Forwarded-For 可能是以","分割的地址列表
		parts := strings.Split(hdrForwardedFor, ",")
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}
		// TODO: 应当返回第一个非本地的地址
		return parts[0]
	}
	return hdrRealIP
}

func logRequestHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {

		uri := r.URL.String()
		ipaddr := requestGetRemoteAddress(r)

		//设置headers
		if len(r.Header) != 0 {
			for k, v := range r.Header {
				w.Header().Add(k, fmt.Sprintf("%v", v))
			}

		}

		//将环境变量写入reponse header
		for _, str := range os.Environ() {
			value := strings.Split(str, "=")
			envkey := value[0]
			envvalue := value[1]
			w.Header().Add(envkey, envvalue)

		}
		m := httpsnoop.CaptureMetrics(h, w, r)
		fmt.Printf("访问的地址：%v,访问IP：%v,返回代码：%d\n", uri, ipaddr, m.Code)
	}
	return http.HandlerFunc(fn)
}
func main() {
	handler := getHealthzHandler{}
	http.Handle("/healthz", logRequestHandler(handler))
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Printf("http server started failed! error: %s\n", err)
	}

}
