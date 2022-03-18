package event

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GetIP(c *gin.Context) string {
	headers := [4]string{
		"HTTP_CF_CONNECTING_IP", "HTTP_X_REAL_IP", "HTTP_X_FORWARDED_FOR", "REMOTE_ADDR",
	}

	for _, headerKey := range headers {
		ip := c.Request.Header.Get(headerKey)
		if ip != "" {
			return ip
		}
	}

	return ""
}

func GetReferer(c *gin.Context) string {
	referer := c.Request.Header.Get("HTTP_REFERER")
	fmt.Println("Referer", referer)
	if referer == "" || referer == "null" {
		return ""
	}
	return referer
}
