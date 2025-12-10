package helpers

import (
	"regexp"
	"strings"
)

func isValidReferrer(referrer string, websiteDomain string) bool {
	isValid := referrer != "" && referrer != "null" && !strings.HasPrefix(referrer, "/") && !strings.HasPrefix(referrer, "#")

	if !isValid {
		return false
	}

	regexDomain := strings.Replace(websiteDomain, ".", `\.`, -1)
	regexString := `^https?:\/\/([a-z0-9-)]+\.)*` + regexDomain
	r := regexp.MustCompile(regexString)
	return !r.MatchString(referrer)
}

func FilterReferrer(ref string, websiteDomain string) (string, string) {
	referrer := strings.ToLower(strings.TrimSpace(ref))
	isValid := isValidReferrer(referrer, websiteDomain)

	if !isValid {
		return "(none)", ""
	}

	domain, fullUrl := CleanupUrl(referrer)

	return domain, fullUrl
}
