package server

import (
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/rs/cors"
)

// NewCORSHandler returns a middleware that handles CORS requests.
// It allows origins from the provided list and origins that share the same base domain.
func NewCORSHandler(allowedOrigins []string, baseDomain string) func(http.Handler) http.Handler {
	log.Println("Allowed origins:", allowedOrigins)
	log.Println("Base domain:", baseDomain)
	c := cors.New(cors.Options{
		AllowOriginFunc: func(origin string) bool {
			log.Printf("CORS: Evaluating origin: %q", origin)
			// 1. Check if origin is in the allowed list
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					log.Printf("CORS: Origin %q allowed by list", origin)
					return true
				}
			}

			// 2. Check if origin shares the same base domain
			if isSameBaseDomain(origin, baseDomain) {
				log.Printf("CORS: Origin %q allowed by base domain %q", origin, baseDomain)
				return true
			}

			log.Printf("CORS: Origin %q NOT allowed (Allowed: %v, BaseDomain: %q)", origin, allowedOrigins, baseDomain)
			return false
		},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders: []string{
			"Origin", "Content-Type", "Accept", "Authorization", "organization_id", "Organization_id", "organization-id", "Organization-Id",
			"X-Grpc-Web", "X-User-Agent", "Grpc-Timeout", "grpc-status", "grpc-message", "grpc-status-details-bin",
			"X-Accept-Content-Transfer-Encoding", "X-Accept-Response-Streaming", "X-Requested-With",
			"Connect-Protocol-Version", "Connect-Timeout-Ms", "Connect-Content-Encoding", "Connect-Accept-Encoding",
		},
		ExposedHeaders: []string{
			"Grpc-Status", "Grpc-Message", "Grpc-Status-Details-Bin", "grpc-status", "grpc-message", "grpc-status-details-bin",
			"X-Grpc-Web", "X-User-Agent", "Connect-Protocol-Version", "organization_id",
		},
		AllowCredentials: true,
		Debug:            true,
	})

	return c.Handler
}

func isSameBaseDomain(origin, baseDomain string) bool {
	log.Println("Origin:", origin)
	if baseDomain == "" {
		return false
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	hostname := u.Hostname()
	log.Println("Origin hostname:", hostname)
	// Check if hostname is exactly the base domain or a subdomain of it
	return hostname == baseDomain || strings.HasSuffix(hostname, "."+baseDomain)
}
