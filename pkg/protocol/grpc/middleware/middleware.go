package middleware

import (
	"bytes"
	"context"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"log"
)

// codeToLevel redirects OK to DEBUG level logging instead of INFO
// This is example how you can log several gRPC code results
func codeToLevel(code codes.Code) zapcore.Level {
	if code == codes.OK {
		// It is DEBUG
		return zap.InfoLevel
	}
	return grpc_zap.DefaultCodeToLevel(code)
}

// AddLogging returns grpc.Server config option that turn on logging.
func AddLogging(logger *zap.Logger, opts []grpc.ServerOption) []grpc.ServerOption {
	// Shared options for the logger, with a custom gRPC code to log level function.
	o := []grpc_zap.Option{
		grpc_zap.WithLevels(codeToLevel),
	}
	// Make sure that log statements internal to gRPC library are logged using the zapLogger as well.
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	// Add unary interceptor
	opts = append(opts, grpc_middleware.WithUnaryServerChain(
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.UnaryServerInterceptor(logger, o...),
	))

	// Add stream interceptor (added as an example here)
	opts = append(opts, grpc_middleware.WithStreamServerChain(
		grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_zap.StreamServerInterceptor(logger, o...),
	))

	return opts
}

//type UnaryServerInterceptor func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error)

//func UnaryServerAuthNInterceptor(opts []grpc.ServerOption) []grpc.ServerOption {
//	log.Println("UnaryServerAuthNInterceptor() called")
//	opts = append(opts, grpc_middleware.WithUnaryServerChain(unaryServerAuthNInterceptor()))
//	return opts
//}

func UnaryServerGrpcCtxTagsInterceptor() grpc.UnaryServerInterceptor {
	return grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor))
}

func UnaryServerAuthNInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log.Println("Authenticated")
		var err error
		//return nil, err
		resp, err := handler(ctx, req)
		return resp, err
	}
}

func UnaryServerLoggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	// Shared options for the logger, with a custom gRPC code to log level function.
	o := []grpc_zap.Option{
		grpc_zap.WithLevels(codeToLevel),
	}
	// Make sure that log statements internal to gRPC library are logged using the zapLogger as well.
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	return grpc_zap.UnaryServerInterceptor(logger, o...)
}

func UnaryServerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	// Shared options for the logger, with a custom gRPC code to log level function.
	o := []grpc_zap.Option{
		grpc_zap.WithLevels(codeToLevel),
	}
	// Make sure that log statements internal to gRPC library are logged using the zapLogger as well.
	grpc_zap.ReplaceGrpcLoggerV2(logger)

	return grpc_zap.UnaryServerInterceptor(logger, o...)
}

func UnaryClientInterceptor(logPayload bool) grpc.UnaryClientInterceptor {
	return func(ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		logRequest(ctx, req, logPayload, "Request sent")
		err := invoker(ctx, method, req, reply, cc, opts...)
		return err
	}
}

func logRequest(ctx context.Context, req interface{}, logPayload bool, format string, args ...interface{}) {
	//entry := For(ctx).log()
	log.Println("Hemant request", req)
	if logPayload {
		if p, ok := req.(proto.Message); ok {
			log.Println("ctx", format, &jsonPbMarshallable{p})
			//entry = entry.Info("", &jsonPbMarshallable{p})
		}
	}
	log.Println("Hemant", format, args)
	//entry.Debugf(format, args...)
}

// Log holds state to log about a request
type Log struct {
	ctx  context.Context
	tags grpc_ctxtags.Tags
}

// For gets a Log for the specified request context. The returned object is a
// valid logrus.FieldLogger.
func For(ctx context.Context) *Log {
	return &Log{ctx, grpc_ctxtags.Extract(ctx)}
}

func (l *Log) log() *zap.Logger {
	return ctxzap.Extract(l.ctx)
}

type jsonPbMarshallable struct {
	proto.Message
}

func (j *jsonPbMarshallable) MarshalJSON() ([]byte, error) {
	b := &bytes.Buffer{}
	if err := JSONPbMarshaller.Marshal(b, j.Message); err != nil {
		return nil, fmt.Errorf("jsonpb serializer failed: %v", err)
	}
	return b.Bytes(), nil
}

var (
	// JSONPbMarshaller is the marshaller used for serializing protobuf messages.
	JSONPbMarshaller = &jsonpb.Marshaler{}
)
