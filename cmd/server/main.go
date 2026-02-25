package main

import (
	"log"
	"net"
	"os"
	"strings"

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

	log.Printf("Admin gRPC server starting on port %s...", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
