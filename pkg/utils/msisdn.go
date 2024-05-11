package utils

import (
	"strings"
)

func SanitizeMsisdn(msisdn string) string {
	switch {
	case strings.HasPrefix(msisdn, "0"):
		return "62" + msisdn[1:]
	case strings.HasPrefix(msisdn, "8"):
		return "62" + msisdn
	case !strings.HasPrefix(msisdn, "8") &&
		!strings.HasPrefix(msisdn, "0") &&
		!strings.HasPrefix(msisdn, "62"):
		return ""
	default:
		return msisdn
	}
}
