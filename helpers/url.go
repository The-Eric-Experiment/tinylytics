package helpers

import (
	"net/url"
	"regexp"
	"strings"
)

func RemoveWWW(url string) string {
	exp := regexp.MustCompile(`^www[\.]`)
	return exp.ReplaceAllString(url, "")
}

func RemoveTrailingSlash(url string) string {
	exp := regexp.MustCompile(`[\/]+$`)
	return exp.ReplaceAllString(url, "")
}

func RemovePrecedingSlash(url string) string {
	exp := regexp.MustCompile(`^\/+`)
	return exp.ReplaceAllString(url, "")
}

func WithQueryPrefix(query string) string {
	if query == "" {
		return ""
	}

	if strings.HasPrefix(query, "?") {
		return query
	}
	return "?" + query
}

func CleanupUrl(input string) (string, string) {
	u, err := url.Parse(input)

	if err != nil {
		return "(none)", ""
	}

	if !strings.HasPrefix(u.Scheme, "http") {
		return "(none)", ""
	}

	if !u.IsAbs() {
		return "(none)", ""
	}

	domain := u.Hostname()
	domain = RemoveWWW(RemoveTrailingSlash(domain))

	path := RemoveTrailingSlash(RemovePrecedingSlash(u.Path))
	query := WithQueryPrefix(RemoveTrailingSlash(RemovePrecedingSlash(u.RawQuery)))
	var fullUrl = ""

	if path != "" {
		if fullUrl == "" {
			fullUrl = domain
		}
		fullUrl += "/" + path
	}

	if query != "" {
		if fullUrl == "" {
			fullUrl = domain
		}
		fullUrl += query
	}

	return domain, fullUrl
}
