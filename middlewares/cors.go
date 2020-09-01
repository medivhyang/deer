package middlewares

import (
	"github.com/medivhyang/deer"
	"net/http"
	"strings"
)

type CORSConfig struct {
	AllowOrigins []string
	AllowMethods []string
	AllowHeaders []string
	ExposeHeaders []string
	AllowCredentials bool
}

func CORS(config ...CORSConfig) deer.Middleware {
	var finalConfig CORSConfig
	if len(config) > 0 {
		finalConfig = config[0]
	} else {
		finalConfig = CORSConfig{}
	}
	return func(h http.Handler) http.Handler {
		return deer.HandlerFunc(func(w *deer.ResponseWriter, r *deer.Request) {
			if len(finalConfig.AllowHeaders) > 0 {
				w.SetHeader("Access-Control-Allow-Origin", strings.Join(finalConfig.AllowHeaders, ","))
			} else {
				w.SetHeader("Access-Control-Allow-Origin", "*")
			}
			if len(finalConfig.AllowMethods) > 0 {
				w.SetHeader("Access-Control-Allow-Methods", strings.Join(finalConfig.AllowMethods, ","))
			} else {
				w.SetHeader("Access-Control-Allow-Methods", "*")
			}
			if len(finalConfig.AllowHeaders) > 0 {
				w.SetHeader("Access-Control-Allow-Headers", strings.Join(finalConfig.AllowHeaders, ","))
			} else {
				w.SetHeader("Access-Control-Allow-Headers", "*")
			}
			if len(finalConfig.ExposeHeaders) > 0 {
				w.SetHeader("Access-Control-Expose-Headers", strings.Join(finalConfig.ExposeHeaders, ","))
			}
			if finalConfig.AllowCredentials {
				w.SetHeader("Access-Control-Allow-Credentials", "true")
			}
			if r.Method() != http.MethodOptions {
				h.ServeHTTP(w.Raw, r.Raw)
				return
			}
		})
	}
}
