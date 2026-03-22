package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	adminpb "persacc/api/v1/admin"
	"persacc/internal/data"
	"persacc/internal/server"

	authpb "github.com/gevorgmb/oauth/api/v1/pb/proto"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	authAddr := os.Getenv("AUTH_SERVICE_ADDR")
	authAddr = strings.TrimSpace(authAddr)
	authAddr = strings.Trim(authAddr, "\"'")
	if authAddr == "" {
		authAddr = "localhost:50061"
		log.Println("AUTH_SERVICE_ADDR not set, using default for client:", authAddr)
	}

	// 1. Initialize Database
	db, err := data.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	log.Println("Database connection established.")

	// 2. Initialize Auth Service Client
	// Use Insecure for now, in prod use TLS
	log.Printf("Dialing auth service at: %s", authAddr)
	authConn, err := grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to auth service: %v", err)
	}
	defer authConn.Close()
	authClient := authpb.NewOAuthClient(authConn)
	log.Printf("Successfully created gRPC client mapped to target %s\n", authAddr)

	// 3. Initialize Admin Server
	srv := server.NewAdminServer(db, authClient)

	// Initialize Auth Interceptor
	authInterceptor := server.NewAuthInterceptor(db, authClient)

	// 4. Start gRPC Server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(authInterceptor.Unary()),
	)
	adminpb.RegisterAdminServiceServer(grpcServer, srv)

	// Register reflection for debugging (grpcurl)
	reflection.Register(grpcServer)

	// 5. Wrap gRPC with gRPC-web
	wrappedGrpc := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool { return true }), // CORS handles origin validation
		grpcweb.WithAllowedRequestHeaders([]string{
			"Origin", "Content-Type", "Accept", "Authorization", "organization_id", "Organization_id", "organization-id", "Organization-Id",
			"X-Grpc-Web", "X-User-Agent", "Grpc-Timeout",
			"Connect-Protocol-Version", "Connect-Timeout-Ms", "Connect-Content-Encoding", "Connect-Accept-Encoding",
		}),
	)

	// CORS configuration
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	allowedOriginsStr = strings.TrimSpace(allowedOriginsStr)
	allowedOriginsStr = strings.Trim(allowedOriginsStr, "\"'")
	var allowedOrigins []string
	if allowedOriginsStr != "" {
		allowedOrigins = strings.Split(allowedOriginsStr, ",")
		for i := range allowedOrigins {
			allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
		}
	}
	baseDomain := os.Getenv("BASE_DOMAIN")
	baseDomain = strings.TrimSpace(baseDomain)
	baseDomain = strings.Trim(baseDomain, "\"'")

	// Create the CORS handler
	corsHandler := server.NewCORSHandler(allowedOrigins, baseDomain)

	// Create the root handler that switches between gRPC-web and standard gRPC
	rootHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Received request: %s %s %s", r.Method, r.URL.Path, r.Proto)
		for name, values := range r.Header {
			for _, value := range values {
				log.Printf("Header: %s: %s", name, value)
			}
		}

		if wrappedGrpc.IsGrpcWebRequest(r) {
			log.Println("Handling as gRPC-web request")
			wrappedGrpc.ServeHTTP(w, r)
			return
		}

		log.Println("Handling as standard gRPC request")
		grpcServer.ServeHTTP(w, r)
	})

	// Wrap rootHandler with CORS
	handlerWithCORS := corsHandler(rootHandler)

	// Wrap everything with h2c as the outermost handler to handle HTTP/2 Cleartext.
	// This ensures that headers (like Origin) are correctly parsed from HTTP/2 streams before being passed to CORS.
	handler := h2c.NewHandler(handlerWithCORS, &http2.Server{})

	log.Printf("Admin server starting on port %s (supporting gRPC, gRPC-web, and CORS)...", port)
	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	log.Printf("Admin server starting on port %s (supporting gRPC, gRPC-web, and CORS)...", port)
	if err := httpServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
