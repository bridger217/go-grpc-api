package grpc

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	firebase "firebase.google.com/go"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	v1 "github.com/bridger217/go-grpc-api/pkg/api/v1"
	"github.com/bridger217/go-grpc-api/pkg/middleware/auth"
)

// RunServer runs gRPC service to publish ToDo service
func RunServer(ctx context.Context, v1API v1.UserServiceServer, port string, fb *firebase.App) error {
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	// setup middleware
	auther := auth.NewAuthenticator(fb)

	// register service
	server := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_logrus.UnaryServerInterceptor(logrus.NewEntry(logrus.New())),
			grpc_auth.UnaryServerInterceptor(auther.Authenticate),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_logrus.StreamServerInterceptor(logrus.NewEntry(logrus.New())),
			grpc_auth.StreamServerInterceptor(auther.Authenticate),
		),
	)
	v1.RegisterUserServiceServer(server, v1API)

	// graceful stop
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		s := <-sigCh
		log.Printf("got signal %v, attempting graceful shutdown", s)
		server.GracefulStop()
		wg.Done()
	}()

	log.Println("starting grpc server on port " + port)
	err = server.Serve(listen)
	if err != nil {
		log.Fatalf("could not serve: %v", err)
	}
	wg.Wait()
	log.Println("clean shutdown")
	return err
}
