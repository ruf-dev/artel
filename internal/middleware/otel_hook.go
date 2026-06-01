package middleware

import (
	"github.com/rs/zerolog"
	otellog "go.opentelemetry.io/otel/log"
	otellogGlobal "go.opentelemetry.io/otel/log/global"
)

type otelLogHook struct {
	logger otellog.Logger
}

func NewOtelLogHook(instrumentationScope string) zerolog.Hook {
	logger := otellogGlobal.GetLoggerProvider().Logger(instrumentationScope)
	return &otelLogHook{logger: logger}
}

func (o *otelLogHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	if e == nil || level == zerolog.NoLevel || level == zerolog.Disabled {
		return
	}
	ctx := e.GetCtx()
	var rec otellog.Record
	rec.SetBody(otellog.StringValue(message))
	rec.SetSeverity(zerologToOtelSeverity(level))
	rec.SetSeverityText(level.String())
	o.logger.Emit(ctx, rec)
}

func zerologToOtelSeverity(level zerolog.Level) otellog.Severity {
	switch level {
	case zerolog.TraceLevel:
		return otellog.SeverityTrace
	case zerolog.DebugLevel:
		return otellog.SeverityDebug
	case zerolog.InfoLevel:
		return otellog.SeverityInfo
	case zerolog.WarnLevel:
		return otellog.SeverityWarn
	case zerolog.ErrorLevel:
		return otellog.SeverityError
	case zerolog.FatalLevel:
		return otellog.SeverityFatal
	case zerolog.PanicLevel:
		return otellog.SeverityFatal4
	default:
		return otellog.SeverityUndefined
	}
}
