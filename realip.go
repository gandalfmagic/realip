package realip

import (
	"context"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/labstack/echo/v4"
)

type contextKey struct{}

const contextKeyString = "github.com/gandalfmagic/realip"

func Get(r *http.Request) string {
	trueClientIP := http.CanonicalHeaderKey("True-Client-IP")
	xRealIP := http.CanonicalHeaderKey("X-Real-IP")
	xForwardedFor := http.CanonicalHeaderKey("X-Forwarded-For")

	var realIP string

	if ip := r.Header.Get(trueClientIP); ip != "" {
		realIP = ip
	} else if ip = r.Header.Get(xRealIP); ip != "" {
		realIP = ip
	} else if ipList := r.Header.Get(xForwardedFor); ipList != "" {
		for _, ip = range strings.Split(ipList, ",") {
			ip = strings.TrimSpace(ip)
			validIp := net.ParseIP(extractIP(ip))
			if validIp != nil && !validIp.IsPrivate() {
				realIP = ip
				break
			}
		}
	}

	if realIP == "" || net.ParseIP(extractIP(realIP)) == nil {
		return ""
	}

	return realIP
}

// Return the ip address from an "ip:port" formatted string
// If the string is already a single ip value, returns that
func extractIP(ip string) string {
	for i := 0; i < len(ip); i++ {
		switch ip[i] {
		case '.':
			if last(ip, ':') < 0 {
				return ip
			}

			ip, _, _ = net.SplitHostPort(ip)
			return ip
		case ':':
			// Normalize IPv6 if needed
			if ip[0] != '[' {
				return ip
			}

			ip, _, _ = net.SplitHostPort(ip)
			return ip
		}
	}

	return ip
}

// Index of rightmost occurrence of b in s.
func last(s string, b byte) int {
	i := len(s)
	for i--; i >= 0; i-- {
		if s[i] == b {
			break
		}
	}
	return i
}

func FromContext(ctx context.Context) string {
	if v, ok := ctx.Value(contextKey{}).(string); ok {
		return v
	}

	return ""
}

func HTTPMiddleware(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if realIP := Get(r); realIP != "" {
			r.WithContext(context.WithValue(r.Context(), contextKey{}, realIP))
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}

func FromEchoContext(c echo.Context) string {
	if v, ok := c.Get(contextKeyString).(string); ok {
		return v
	}

	return ""
}

func EchoMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if realIP := Get(c.Request()); realIP != "" {
			c.Set(contextKeyString, realIP)
		}

		return next(c)
	}
}

func FromGinContext(c *gin.Context) string {
	if v, ok := c.Value(contextKeyString).(string); ok {
		return v
	}

	return ""
}

func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if realIP := Get(c.Request); realIP != "" {
			c.Set(contextKeyString, realIP)
		}

		c.Next()
	}
}
