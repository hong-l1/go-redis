package logger

import "go.uber.org/zap"

type ZapLogger struct {
	l *zap.Logger
}

func (z *ZapLogger) With(args ...Field) Loggerv1 {
	return &ZapLogger{
		l: z.l.With(z.toZapFeild(args)...),
	}
}

func NewZapLogger(l *zap.Logger) Loggerv1 {
	return &ZapLogger{l: l}
}
func (z *ZapLogger) Debug(msg string, args ...Field) {
	z.l.Debug(msg, z.toZapFeild(args)...)
}

func (z *ZapLogger) Info(msg string, args ...Field) {
	z.l.Info(msg, z.toZapFeild(args)...)
}

func (z *ZapLogger) Warn(msg string, args ...Field) {
	z.l.Warn(msg, z.toZapFeild(args)...)
}

func (z *ZapLogger) Error(msg string, args ...Field) {
	z.l.Error(msg, z.toZapFeild(args)...)
}
func (z *ZapLogger) toZapFeild(args []Field) []zap.Field {
	ans := make([]zap.Field, 0, len(args))
	for _, arg := range args {
		ans = append(ans, zap.Any(arg.Key, arg.Value))
	}
	return ans
}
