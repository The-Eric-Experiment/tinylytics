package event

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func GetIP(c *gin.Context) string {
	headers := []string{
		"CF-Connecting-IP",
		"X-Real-IP",
		"X-Forwarded-For",
		"HTTP_CF_CONNECTING_IP",
		"HTTP_X_REAL_IP",
		"HTTP_X_FORWARDED_FOR",
		"REMOTE_ADDR",
	}

	for _, headerKey := range headers {
		ip := c.Request.Header.Get(headerKey)
		if ip != "" {
			// Handle comma-separated lists (take first IP)
			if strings.Contains(ip, ",") {
				parts := strings.Split(ip, ",")
				return strings.TrimSpace(parts[0])
			}
			return ip
		}
	}

	// Fallback to connection IP
	return c.ClientIP()
}

func GetReferer(c *gin.Context) string {
	headers := []string{
		"Referer",      // Standard HTTP header
		"HTTP_REFERER", // CGI-style
	}

	for _, headerKey := range headers {
		referer := c.Request.Header.Get(headerKey)
		if referer != "" && referer != "null" {
			return referer
		}
	}

	return ""
}
