package postgres

import (
	"context"

	"github.com/altinity/transfer/pkg/contextutil"
	"github.com/altinity/transfer/pkg/util"
	"github.com/jackc/pgx/v4"
	"go.ytsaurus.tech/library/go/core/log"
)

type pgxLogger struct {
	logger log.Logger
}

var pgLoggerNotToLog = contextutil.NewContextKey()

func withNotToLog(ctx context.Context) context.Context {
	return context.WithValue(ctx, pgLoggerNotToLog, true)
}

func isNotToLogInContext(ctx context.Context) bool {
	_, ok := ctx.Value(pgLoggerNotToLog).(bool)
	return ok
}

func (p pgxLogger) Log(ctx context.Context, level pgx.LogLevel, msg string, data map[string]interface{}) {
	if isNotToLogInContext(ctx) {
		return
	}

	var params []log.Field
	for k, v := range data {
		if k == "sql" {
			// https://github.com/altinity/transfer/arcadia/vendor/github.com/jackc/pgx/v4/conn.go?rev=r9171541#L413
			query, ok := v.(string)
			if ok {
				v = util.DefaultSample(query)
			}
		}
		params = append(params, log.Any(k, v))
	}
	switch level {
	case pgx.LogLevelInfo:
		// Info level is too verbose, hide it behind debug
		p.logger.Debug(msg, params...)
	case pgx.LogLevelWarn:
		p.logger.Warn(msg, params...)
	case pgx.LogLevelError:
		params = append(params, log.Any("callstack", util.GetCurrentGoroutineCallstack()))
		p.logger.Error(msg, params...)
	case pgx.LogLevelDebug:
		p.logger.Debug(msg, params...)
	}
}

func WithLogger(connConfig *pgx.ConnConfig, lgr log.Logger) *pgx.ConnConfig {
	connConfig.Logger = pgxLogger{logger: lgr}
	return connConfig
}
