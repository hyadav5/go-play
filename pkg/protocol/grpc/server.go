package grpc

import (
	"context"
	"github.com/go-play/pkg/api/v1/v1"
	"github.com/go-play/pkg/logger"
	"github.com/go-play/pkg/protocol/grpc/middleware"
	//mlog "github.com/go-play/pkg/protocol/grpc/middleware/logging"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
)

// RunServer runs gRPC service to publish ToDo service
func RunServer(ctx context.Context, v1API v1.ToDoServiceServer, port string) error {
	log.Println("Running gRPC server. RunServer(). Port: ", port, ", ctx: ", ctx)
	listen, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	//opts := []grpc.ServerOption{}

	// add middleware
	//opts = middleware.AddLogging(logger.Log, opts)

	var authEnabled = true
	var middlewareLogging = true
	var unaryInterceptors []grpc.UnaryServerInterceptor
	// var unaryClientInterceptors []grpc.UnaryClientInterceptor

	// Some default interceptors
	// `grpc_ctxtags` adds a Tag object to the context that can be used by other middleware to add context about a request
	unaryInterceptors = append(unaryInterceptors, middleware.UnaryServerGrpcCtxTagsInterceptor())

	if authEnabled {
		unaryInterceptors = append(unaryInterceptors, middleware.UnaryServerAuthNInterceptor())
	}
	if middlewareLogging {
		unaryInterceptors = append(unaryInterceptors, middleware.UnaryServerInterceptor(logger.Log))
		//unaryClientInterceptors = append(unaryClientInterceptors, middleware.UnaryClientInterceptor(true))
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(unaryInterceptors...)),
	}

	// register service
	server := grpc.NewServer(opts...)
	v1.RegisterToDoServiceServer(server, v1API)

	// graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for range c {
			// sig is a ^C, handle it
			log.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	// start gRPC server
	log.Println("starting gRPC server...")
	return server.Serve(listen)
}
