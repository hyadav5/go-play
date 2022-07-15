package middleware

import (
	"bytes"
	"fmt"
	"google.golang.org/grpc"
)

/*
Copyright (c) 2019 VMware, Inc. All rights reserved.

Proprietary and confidential.

Unauthorized copying or use of this file, in any medium or form,
is strictly prohibited.
*/

import (
	"context"
	//nolint:staticcheck
	"github.com/golang/protobuf/jsonpb"
	//nolint:staticcheck
	"github.com/golang/protobuf/proto"
)

// Source: https://github.com/grpc-ecosystem/go-grpc-middleware/blob/master/logging/logrus/payload_interceptors.go

// NOTE: Above payload interceptor logs request and response both. But we don't want to log response (Response will be big in our case).
// So, we took above implementation and removed response logging. This will also give us more control.
// E.g. If we want to log stack trace in case of some error codes later.

var (
	// JSONPbMarshaller is the marshaller used for serializing protobuf messages.
	JSONPbMarshaller = &jsonpb.Marshaler{}
)

// UnaryServerInterceptor returns a new unary server interceptor that logs a time marker for every request received.
// If the logPayload flag is set, it also logs the payloads of requests.
//
// This *only* works when placed *after* the `grpc_logrus.UnaryServerInterceptor`. However, the logging can be done to a
// separate instance of the logger.
func UnaryServerInterceptor(logPayload bool) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logRequest(ctx, req, logPayload, "Request received")
		resp, err := handler(ctx, req)
		return resp, err
	}
}

// StreamServerInterceptor returns a new unary server interceptor that logs a time marker for every message received.
// If the logPayload flag is set, it also logs the payloads of message.
//
// This *only* works when placed *after* the `grpc_logrus.StreamServerInterceptor`. However, the logging can be done to a
// separate instance of the logger.
func StreamServerInterceptor(logPayload bool) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		newStream := &loggingServerStream{ServerStream: stream, logPayload: logPayload}
		return handler(srv, newStream)
	}
}

type loggingServerStream struct {
	grpc.ServerStream
	logPayload bool
}

func (l *loggingServerStream) SendMsg(m interface{}) error {
	return l.ServerStream.SendMsg(m)
}

func (l *loggingServerStream) RecvMsg(m interface{}) error {
	err := l.ServerStream.RecvMsg(m)
	if err == nil {
		logRequest(l.Context(), m, l.logPayload, "Message received")
	}
	return err
}

func logRequest(ctx context.Context, req interface{}, logPayload bool, format string, args ...interface{}) {
	entry := For(ctx).log()
	if logPayload {
		if _, ok := req.(proto.Message); ok {
			//entry.With(zap.Interface("request.msg", string(jsonPbMarshallable{p})))
			//entry = entry.With(zap.String("request.msg", )"request.msg", &jsonPbMarshallable{p})
		}
	}
	entry.Sugar().Debugf(format, args...)
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
