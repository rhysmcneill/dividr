package middleware

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rhysmcneill/dividr/internal/hmrc"
)

// Context Key to retrieve metadata later
type fraudContextKey string

const FraudMetadataKey fraudContextKey = "hmrc_fraud_metadata"

// FraudPreventionMiddleware captures necessary headers for HMRC compliance
func FraudPreventionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Handle Device ID (Persistent Cookie)
		deviceID := ""
		cookie, err := r.Cookie("hmrc_device_id")
		if err == nil {
			deviceID = cookie.Value
		} else {
			// Generate new Device ID
			deviceID = uuid.New().String()
			// Set a long-lived cookie (18 months is standard for HMRC persistence)
			http.SetCookie(w, &http.Cookie{
				Name:     "hmrc_device_id",
				Value:    deviceID,
				Path:     "/",
				Expires:  time.Now().Add(540 * 24 * time.Hour), // ~18 months
				HttpOnly: true,
				Secure:   true, // Ensure you are running HTTPS in prod
				SameSite: http.SameSiteLaxMode,
			})
		}

		// 2. Capture IP (Handling Reverse Proxies like Cloudflare/Nginx)
		ip, port, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
			port = ""
		}

		// Priority: X-Forwarded-For -> CF-Connecting-IP -> RemoteAddr
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			// XFF can be "client, proxy1, proxy2" - take the first one
			parts := strings.Split(forwarded, ",")
			ip = strings.TrimSpace(parts[0])
		}

		// 3. Construct Metadata
		meta := hmrc.FraudMetadata{
			ClientIP:        ip,
			ClientPort:      port,
			ClientDeviceID:  deviceID,
			ClientUserAgent: r.UserAgent(),
			ClientTimezone:  "UTC",
			// ClientUserID will be filled in by the handler after Auth check
		}

		// 4. Store in Context
		ctx := context.WithValue(r.Context(), FraudMetadataKey, meta)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetFraudMetadata helper function
func GetFraudMetadata(ctx context.Context) (hmrc.FraudMetadata, bool) {
	meta, ok := ctx.Value(FraudMetadataKey).(hmrc.FraudMetadata)
	return meta, ok
}
