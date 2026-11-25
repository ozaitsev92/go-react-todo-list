package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// CORS allows credentialed requests from allowedOrigin only; other origins'
// preflights get 403.
func CORS(allowedOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.Writer.Header()
		h.Add("Vary", "Origin")

		origin := c.GetHeader("Origin")
		allowed := origin != "" && origin == allowedOrigin

		if allowed {
			h.Set("Access-Control-Allow-Origin", origin)
			h.Set("Access-Control-Allow-Credentials", "true")
		}

		isPreflight := c.Request.Method == http.MethodOptions &&
			origin != "" && c.GetHeader("Access-Control-Request-Method") != ""
		if !isPreflight {
			c.Next()
			return
		}

		if !allowed {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}

		h.Add("Vary", "Access-Control-Request-Method")
		h.Add("Vary", "Access-Control-Request-Headers")
		h.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		h.Set("Access-Control-Allow-Headers", "Content-Type, X-Request-ID")
		h.Set("Access-Control-Max-Age", "600")
		c.AbortWithStatus(http.StatusNoContent)
	}
}
