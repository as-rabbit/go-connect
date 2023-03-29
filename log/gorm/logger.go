package logger

import (
	"context"
	"database/sql/driver"
	"errors"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	opentracinglog "github.com/opentracing/opentracing-go/log"
	"github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	"gorm.io/gorm"
	glogger "gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
	"time"
)

type dbLogger struct {
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
	Config
}

type Config struct {
	SlowThreshold             time.Duration
	IgnoreRecordNotFoundError bool
	LogLevel                  glogger.LogLevel
	Dsn                       string
}

func NewGORMLogger(config Config) glogger.Interface {
	var (
		infoStr      = "%s [info] "
		warnStr      = "%s [warn] "
		errStr       = "%s [error] "
		traceStr     = "%s [%.3fms] [rows:%v] %s"
		traceWarnStr = "%s %s [%.3fms] [rows:%v] %s"
		traceErrStr  = "%s %s [%.3fms] [rows:%v] %s"
	)

	return &dbLogger{
		Config:       config,
		infoStr:      infoStr,
		warnStr:      warnStr,
		errStr:       errStr,
		traceStr:     traceStr,
		traceWarnStr: traceWarnStr,
		traceErrStr:  traceErrStr,
	}
}

func (log *dbLogger) LogMode(level glogger.LogLevel) glogger.Interface {
	l := *log
	l.LogLevel = level
	return &l
}

func (log *dbLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if log.LogLevel >= glogger.Info {
		logrus.Infof(log.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, args...)...)
	}
}

func (log *dbLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if log.LogLevel >= glogger.Warn {
		logrus.Warnf(log.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, args...)...)
	}
}

func (log *dbLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if log.LogLevel >= glogger.Error {
		logrus.Errorf(log.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, args...)...)
	}
}

// Trace print sql message
func (log *dbLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if log.LogLevel <= glogger.Silent {
		return
	}

	//newCtx, ok := tls.GetContext()
	//if ok {
	//	ctx = newCtx
	//}

	elapsed := time.Since(begin)
	switch {
	case err != nil && log.LogLevel >= glogger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !log.IgnoreRecordNotFoundError):
		sql, rows := fc()
		log.tryTrace(ctx, sql, begin, err)
		//logW := logger.For(ctx)
		if rows == -1 {
			//logW.Errorf(log.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			//logW.Errorf(log.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	case elapsed > log.SlowThreshold && log.SlowThreshold != 0:
		sql, rows := fc()
		//logW := logger.DefaultKit.Slow().For(ctx)
		log.tryTrace(ctx, sql, begin, err)
		//slowLog := fmt.Sprintf("slow sql >= %v", log.SlowThreshold)
		if rows == -1 {
			//logW.Errorf(log.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			//logW.Errorf(log.traceWarnStr, utils.FileWithLineNum(), slowLog, float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	default:
		sql, rows := fc()
		//logW := logger.For(ctx)
		log.tryTrace(ctx, sql, begin, err)
		if rows == -1 {
			//logW.Debugf(log.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, "-", sql)
		} else {
			//logW.Debugf(log.traceStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
		}
	}
}

func (log *dbLogger) tryTrace(ctx context.Context, query string, startTime time.Time, err error) string {
	var (
		traceID string
	)

	if err == driver.ErrSkip {
		return ""
	}

	span, _ := opentracing.StartSpanFromContext(ctx, "MYSQL Client", opentracing.StartTime(startTime))
	// annotation
	ext.PeerAddress.Set(span, log.Dsn)
	ext.DBType.Set(span, "database/mysql")
	ext.Component.Set(span, "golang/mysql-client")

	//ext.DBStatement.Set(span, query)
	if sc, ok := span.Context().(jaeger.SpanContext); ok {
		traceID = sc.TraceID().String()
	}

	span.LogFields(opentracinglog.String("query", query))
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(opentracinglog.Error(err))
	}
	span.Finish()
	return traceID
}
