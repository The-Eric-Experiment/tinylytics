package routes

import (
	"encoding/json"
	"net/http"
	"time"
	"tinylytics/event"

	"github.com/gin-gonic/gin"
)

func PostEvent(eventQueue *event.EventQueue) func(c *gin.Context) {
	return func(c *gin.Context) {
		uagent := c.Request.Header.Get("User-Agent")
		chUa := c.Request.Header.Get("Sec-CH-UA")
		chMobile := c.Request.Header.Get("Sec-CH-UA-Mobile")
		chPlatform := c.Request.Header.Get("Sec-CH-UA-Platform")
		chVersion := c.Request.Header.Get("Sec-CH-UA-Full-Version")
		chPlatformVersion := c.Request.Header.Get("Sec-CH-UA-Platform-Version")

		if chUa == "" && chMobile == "" && chPlatform == "" && chVersion == "" && chPlatformVersion == "" {
			c.Header("Accept-CH", "sec-ch-ua,sec-ch-ua-platform,sec-ch-ua-mobile,sec-ch-ua-full-version,Sec-CH-UA-Platform-Version,sec-ch-width,width,sec-ch-viewport-width,viewport-width")
		}

		var ed event.EventData
		err := json.NewDecoder(c.Request.Body).Decode(&ed)

		if err != nil {
			c.String(http.StatusBadRequest, "There's an issue with the event data")
			return
		}

		if ed.Name != "pageview" {
			c.String(http.StatusBadRequest, "Only the 'pageview' event is supported at the moment")
			return
		}

		if ed.Domain == "" {
			c.String(http.StatusBadRequest, "No domain was set")
			return
		}

		if ed.Page == "" {
			c.String(http.StatusBadRequest, "No page was set")
			return
		}

		info := &event.ClientInfo{
			Name:                      ed.Name,
			UserAgent:                 uagent,
			IP:                        event.GetIP(c),
			HostName:                  c.Request.Host,
			Domain:                    ed.Domain,
			Page:                      ed.Page,
			ClientHintUA:              chUa,
			ClientHintMobile:          chMobile,
			ClientHintPlatform:        chPlatform,
			ClientHintFullVersion:     chVersion,
			ClientHintPlatformVersion: chPlatformVersion,
			Referer:                   event.GetReferer(c),
			Time:                      time.Now().UTC(),
			ScreenWidth:               ed.ScreenWidth,
		}

		eventQueue.Push(info)

		c.String(http.StatusOK, "ok")
	}
}
