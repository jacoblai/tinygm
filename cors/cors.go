package cors

import (
	"net/http"
)

const (
	options      = "OPTIONS"
	allowOrigin  = "Access-Control-Allow-Origin"
	allowMethods = "Access-Control-Allow-Methods"
	allowHeaders = "Access-Control-Allow-Headers"
	origin       = "Origin"
	methods      = "GET,PUT,POST,DELETE,PATCH"
	// If you want to expose some other headers add it here
	headers = "Authorization,Content-Length,Content-Type"
)

// Handler will allow cross-origin HTTP requests
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set allow origin to match origin of our request or fall back to *
		//if o := r.Header.Get(origin); o != "" {
		//	w.Header().Set(allow_origin, o)
		//} else {
		//if engine.DEBUG {
		//switch r.Method {
		//case "POST", "GET", "PUT", "DELETE", "HEAD":
		//	log.Printf("new request from IP:%v;request METHOD %v;Host %v;request path %v", r.RemoteAddr, r.Method, r.Host, r.URL.Path)
		//}
		//}
		w.Header().Set(allowOrigin, "*")
		//}

		//if o := r.Header.Get(origin); o != "" {
		//	dm, err := url.Parse(o)
		//	if err == nil {
		//		parts := strings.Split(dm.Hostname(), ".")
		//		if len(parts) >= 2 {
		//			domain := parts[len(parts)-2] + "." + parts[len(parts)-1]
		//			if len(parts) > 2 && domain == "net.cn" {
		//				domain = parts[len(parts)-3] + "." + domain
		//			}
		//			if domain == "caredaily.com" || domain == "hkaspire.net.cn" || domain == "hkaspire.cn" || strings.HasPrefix(dm.Hostname(), "192.168") {
		//				w.Header().Set(allow_origin, o)
		//			}
		//		}
		//	}
		//}

		// Set other headers
		w.Header().Set(allowHeaders, headers)
		w.Header().Set(allowMethods, methods)

		// If this was preflight options request let's write empty ok response and return
		if r.Method == options {
			w.WriteHeader(http.StatusOK)
			w.Write(nil)
			return
		}
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		//start := time.Now()
		next.ServeHTTP(w, r)
		//log.Println(time.Now().Sub(start), r.URL.Path)
	})
}
