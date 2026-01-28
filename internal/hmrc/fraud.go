package hmrc

import (
	"net/http"
	"net/url"
)

// FraudMetadata holds the mandatory anti-fraud data
type FraudMetadata struct {
	ClientIP        string
	ClientPort      string
	ClientDeviceID  string
	ClientUserAgent string
	ClientTimezone  string
	ClientUserID    string
}

func AddFraudHeaders(req *http.Request, meta FraudMetadata) {
	req.Header.Set("Gov-Client-Connection-Method", "WEB_APP_VIA_SERVER")

	if meta.ClientIP != "" {
		req.Header.Set("Gov-Client-Public-IP", meta.ClientIP)
	}
	if meta.ClientPort != "" {
		req.Header.Set("Gov-Client-Public-Port", meta.ClientPort)
	}
	if meta.ClientDeviceID != "" {
		req.Header.Set("Gov-Client-Device-ID", meta.ClientDeviceID)
	}

	// We must URL-encode it just in case, though UUIDs are safe.
	if meta.ClientUserID != "" {
		req.Header.Set("Gov-Client-User-IDs", "dividr="+url.QueryEscape(meta.ClientUserID))
	}

	req.Header.Set("Gov-Client-Timezone", "UTC")

	if meta.ClientUserAgent != "" {
		req.Header.Set("Gov-Client-User-Agent", meta.ClientUserAgent)
	}

	// NOTE: I am omitting Gov-Client-Screens and Gov-Client-Window-Size for now.
	// If Sandbox rejects this, we will add a JavaScript collector later.
}
