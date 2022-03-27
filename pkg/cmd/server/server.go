package cmd

import (
	"context"
	"database/sql"
	"encoding/base64"
	"flag"
	"fmt"

	firebase "firebase.google.com/go"
	"google.golang.org/api/option"

	"github.com/bridger217/go-grpc-api/pkg/env"
	"github.com/bridger217/go-grpc-api/pkg/protocol/grpc"
	"github.com/bridger217/go-grpc-api/pkg/protocol/rest"
	v1 "github.com/bridger217/go-grpc-api/pkg/service/v1"
)

// Config is configuration for Server
type Config struct {
	// gRPC server start parameters section
	// gRPC is TCP port to listen by gRPC server
	GRPCPort string

	// HTTP/REST gateway start parameters section
	// HTTPPort is TCP port to listen by HTTP/REST gateway
	HTTPPort string
}

// RunServer runs gRPC server and HTTP gateway
func RunServer() error {
	ctx := context.Background()

	var cfg Config
	flag.StringVar(&cfg.GRPCPort, "grpc-port", "9090", "gRPC port to bind")
	flag.StringVar(&cfg.HTTPPort, "http-port", "8080", "HTTP port to bind")
	flag.Parse()

	if len(cfg.GRPCPort) == 0 {
		return fmt.Errorf("invalid TCP port for gRPC server: '%s'", cfg.GRPCPort)
	}
	if len(cfg.HTTPPort) == 0 {
		return fmt.Errorf("invalid TCP port for HTTP gateway: '%s'", cfg.HTTPPort)
	}

	var eMan env.EnvManager

	// DB setup
	dbStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
		eMan.GetDbUser(),
		eMan.GetDbPassword(),
		eMan.GetDbIpAddr(),
		eMan.GetDbName())
	db, err := sql.Open("mysql", dbStr)
	defer db.Close()
	if err != nil {
		return fmt.Errorf("mysql setup failed: '%s", err)
	}

	// Firebase setup
	fbJSON, err := base64.StdEncoding.DecodeString(eMan.GetFirebaseJson())
	if err != nil {
		return fmt.Errorf("firebase credential decode failed: '%s", err)
	}
	fb, err := firebase.NewApp(
		context.Background(),
		nil,
		option.WithCredentialsJSON([]byte(fbJSON)))
	if err != nil {
		return fmt.Errorf("firebase setup failed: '%s", err)
	}

	v1API := v1.NewUserServiceServer(db)
	// run HTTP gateway
	go func() {
		_ = rest.RunServer(ctx, cfg.GRPCPort, cfg.HTTPPort)
	}()

	return grpc.RunServer(ctx, v1API, cfg.GRPCPort, fb)
}
