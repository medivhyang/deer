package middlewares

import (
	"github.com/medivhyang/deer"
	"strings"
)

type CORSConfig struct {
	AllowOrigins     []string
	AllowMethods     []string
	AllowHeaders     []string
	ExposeHeaders    []string
	AllowCredentials bool
}

func CORS(config ...CORSConfig) deer.Middleware {
	var finalConfig CORSConfig
	if len(config) > 0 {
		finalConfig = config[0]
	} else {
		finalConfig = CORSConfig{
			AllowOrigins: []string{"*"},
			AllowMethods: []string{"*"},
			AllowHeaders: []string{"*"},
		}
	}
	return func(h deer.HandlerFunc) deer.HandlerFunc {
		return func(w deer.ResponseWriter, r *deer.Request) {
			if len(finalConfig.AllowOrigins) > 0 {
				w.Header("Access-Control-Allow-Origin", strings.Join(finalConfig.AllowOrigins, ","))
			}
			if len(finalConfig.AllowMethods) > 0 {
				w.Header("Access-Control-Allow-Methods", strings.Join(finalConfig.AllowMethods, ","))
			}
			if len(finalConfig.AllowHeaders) > 0 {
				w.Header("Access-Control-Allow-Headers", strings.Join(finalConfig.AllowHeaders, ","))
			}
			if len(finalConfig.ExposeHeaders) > 0 {
				w.Header("Access-Control-Expose-Headers", strings.Join(finalConfig.ExposeHeaders, ","))
			}
			if finalConfig.AllowCredentials {
				w.Header("Access-Control-Allow-Credentials", "true")
			}
			h.Next(w, r)
		}
	}
}
