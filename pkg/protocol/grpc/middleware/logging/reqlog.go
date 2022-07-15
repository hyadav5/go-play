package middleware

/*
Copyright (c) 2019 VMware, Inc. All rights reserved.

Proprietary and confidential.

Unauthorized copying or use of this file, in any medium or form,
is strictly prohibited.
*/

// Package reqlog contains helpers for using per-request logging

import (
	"context"
	"go.uber.org/zap"
	//"github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus/ctxlogrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
)

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

// Tag adds a new tag to the request context. This is currently a thin wrapper
// around grpc_ctxtags.Extract(ctx).Set().
func (l *Log) Tag(field string, value interface{}) {
	l.tags.Set(field, value)
}

func (l *Log) log() *zap.Logger {
	return ctxzap.Extract(l.ctx)
}

// Everything below this point is just boring wrappers around *logrus.Entry
// to make *Log a valid logrus.FieldLogger. We have to wrap all of these
// rather than just having our Log struct embed a logrus.Entry because we need
// to get a fresh logrus.Entry from the context for every call, in case
// the user has set new tags.

// WithField creates an entry from the request logger and adds a field to
// it. If you want multiple fields, use `WithFields`.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
//func (l *Log) WithField(key string, value interface{}) *logrus.Entry {
//	return l.log().With()
//}

// WithFields creates an entry from the request logger and adds multiple
// fields to it. This is simply a helper for `WithField`, invoking it
// once for each field.
//
// Note that it doesn't log until you call Debug, Print, Info, Warn, Fatal
// or Panic on the Entry it returns.
//func (l *Log) WithFields(fields logrus.Fields) *logrus.Entry {
//	return l.log().WithFields(fields)
//}

// WithError creates an entry from the request logger and adds an error to it,
// using the value defined in ErrorKey as key.
//func (l *Log) WithError(err error) *logrus.Entry {
//	return l.log().WithError(err)
//}

//// Debug logs a message at level Debug on the standard logger.
//func (l *Log) Debug(args ...interface{}) { l.log().Debug(args...) }
//
//// Info logs a message at level Info on the standard logger.
//func (l *Log) Info(args ...interface{}) { l.log().Info(args...) }
//
//// Print logs a message at level Info on the standard logger.
//func (l *Log) Print(args ...interface{}) { l.log().Print(args...) }
//
//// Warn logs a message at level Warn on the standard logger.
//func (l *Log) Warn(args ...interface{}) { l.log().Warn(args...) }
//
//// Warning logs a message at level Warn on the standard logger.
//func (l *Log) Warning(args ...interface{}) { l.log().Warning(args...) }
//
//// Error logs a message at level Error on the standard logger.
//func (l *Log) Error(args ...interface{}) { l.log().Error(args...) }
//
//// Panic logs a message at level Panic on the standard logger.
//func (l *Log) Panic(args ...interface{}) { l.log().Panic(args...) }
//
//// Fatal logs a message at level Fatal on the standard logger.
//func (l *Log) Fatal(args ...interface{}) { l.log().Fatal(args...) }
//
//// Debugf logs a message at level Debug on the standard logger.
//func (l *Log) Debugf(format string, args ...interface{}) { l.log().Debugf(format, args...) }
//
//// Infof logs a message at level Info on the standard logger.
//func (l *Log) Infof(format string, args ...interface{}) { l.log().Infof(format, args...) }
//
//// Printf logs a message at level Info on the standard logger.
//func (l *Log) Printf(format string, args ...interface{}) { l.log().Printf(format, args...) }
//
//// Warnf logs a message at level Warn on the standard logger.
//func (l *Log) Warnf(format string, args ...interface{}) { l.log().Warnf(format, args...) }
//
//// Warningf logs a message at level Warn on the standard logger.
//func (l *Log) Warningf(format string, args ...interface{}) { l.log().Warningf(format, args...) }
//
//// Errorf logs a message at level Error on the standard logger.
//func (l *Log) Errorf(format string, args ...interface{}) { l.log().Errorf(format, args...) }
//
//// Panicf logs a message at level Panic on the standard logger.
//func (l *Log) Panicf(format string, args ...interface{}) { l.log().Panicf(format, args...) }
//
//// Fatalf logs a message at level Fatal on the standard logger.
//func (l *Log) Fatalf(format string, args ...interface{}) { l.log().Fatalf(format, args...) }
//
//// Debugln logs a message at level Debug on the standard logger.
//func (l *Log) Debugln(args ...interface{}) { l.log().Debugln(args...) }
//
//// Infoln logs a message at level Info on the standard logger.
//func (l *Log) Infoln(args ...interface{}) { l.log().Infoln(args...) }
//
//// Println logs a message at level Info on the standard logger.
//func (l *Log) Println(args ...interface{}) { l.log().Println(args...) }
//
//// Warnln logs a message at level Warn on the standard logger.
//func (l *Log) Warnln(args ...interface{}) { l.log().Warnln(args...) }
//
//// Warningln logs a message at level Warn on the standard logger.
//func (l *Log) Warningln(args ...interface{}) { l.log().Warningln(args...) }
//
//// Errorln logs a message at level Error on the standard logger.
//func (l *Log) Errorln(args ...interface{}) { l.log().Errorln(args...) }
//
//// Panicln logs a message at level Panic on the standard logger.
//func (l *Log) Panicln(args ...interface{}) { l.log().Panicln(args...) }
//
//// Fatalln logs a message at level Fatal on the standard logger.
//func (l *Log) Fatalln(args ...interface{}) { l.log().Fatalln(args...) }
