package log

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugar *zap.SugaredLogger

func newZapEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:       "message",                      // 输入信息的key名
		LevelKey:         "level",                        // 输出日志级别的key名
		TimeKey:          "time",                         // 输出时间的key名
		NameKey:          "logger",                       // 日志信息key名
		CallerKey:        "caller",                       // 调用者key名
		StacktraceKey:    "stacktrace",                   // 调用栈key名
		LineEnding:       zapcore.DefaultLineEnding,      // 每行分隔符,默认\n
		EncodeLevel:      CustomLevelEncoder,             // level值的封装,配置为序列化为全大写
		EncodeTime:       TimeEncoder,                    // 时间格式,配置为[2006-01-02 15:04:05]
		EncodeDuration:   zapcore.SecondsDurationEncoder, // 执行消耗时间格式,配置为浮点秒
		EncodeCaller:     zapcore.ShortCallerEncoder,     // 调用者格式,配置为包/文件:行号
		EncodeName:       zapcore.FullNameEncoder,        // 日志信息名处理,默认无处理
		ConsoleSeparator: " ",                            // 日志的分离器，以什么符号分割
	}
}

func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + level.CapitalString() + "]")
}

func TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.999999"))
}

func getLogLevel(logLevel string) zapcore.Level {
	var zapLevel zapcore.Level
	switch logLevel {
	case "debug":
		zapLevel = zap.DebugLevel
	case "info":
		zapLevel = zap.InfoLevel
	case "warn":
		zapLevel = zap.WarnLevel
	case "error":
		zapLevel = zap.ErrorLevel
	default:
		zapLevel = zap.InfoLevel
	}
	return zapLevel
}

func New(logLevel string) {
	sugar = NewLogger(logLevel)
	Info("log level is", logLevel)
}

func NewLogger(logLevel string) *zap.SugaredLogger {
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(newZapEncoderConfig()),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout)),
		getLogLevel(logLevel),
	)

	logger := zap.New(
		core,
		zap.AddCaller(),                       // 日志增加调用者
		zap.AddStacktrace(zapcore.PanicLevel), // panic级别下增加调用栈
		zap.AddCallerSkip(1),                  // 如果日志方法有封装,则调用者输出跳过的层数
	)
	sugar := logger.Sugar()
	return sugar
}

func Debug(args ...interface{}) {
	sugar.Debug(args)
}

func Info(args ...interface{}) {
	sugar.Info(args)
}

func Warn(args ...interface{}) {
	sugar.Warn(args)
}

func Error(args ...interface{}) {
	sugar.Error(args)
}

func Panic(args ...interface{}) {
	sugar.Panic(args)
}

func Fatal(args ...interface{}) {
	sugar.Fatal(args)
}

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args)
}

func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args)
}

func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args)
}

func Errorf(template string, args ...interface{}) {
	sugar.Error(template, args)
}

func Panicf(template string, args ...interface{}) {
	sugar.Panicf(template, args)
}

func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args)
}
