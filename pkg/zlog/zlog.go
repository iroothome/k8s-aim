package zlog

import (
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Options struct {
	LogFileDir    string //文件保存地方
	AppName       string //日志文件前缀
	ErrorFileName string
	WarnFileName  string
	InfoFileName  string
	DebugFileName string
	Level         zapcore.Level //日志等级
	MaxSize       int           //日志文件小大（M）
	MaxBackups    int           // 最多存在多少个切片文件
	MaxAge        int           //保存的最大天数
	Development   bool          //是否是开发模式
	zap.Config
}

type ModOptions func(options *Options)

var (
	zLogger                        *zap.SugaredLogger
	l                              *Logger
	sp                             = string(filepath.Separator)
	errWS, warnWS, infoWS, debugWS zapcore.WriteSyncer       // IO输出
	debugConsoleWS                 = zapcore.Lock(os.Stdout) // 控制台标准输出
	errorConsoleWS                 = zapcore.Lock(os.Stderr)
)

type Logger struct {
	*zap.Logger
	sync.RWMutex
	Opts      *Options `json:"opts"`
	zapConfig zap.Config
	inited    bool
}

func init() {
	lg := NewLogger(SetAppName("k8s-aim"), SetDevelopment(true), SetLevel(zap.DebugLevel), SetErrorFileName("error.log"))
	zLogger = lg.Sugar()
}

func NewLogger(mod ...ModOptions) *zap.Logger {
	l = &Logger{}
	l.Lock()
	defer l.Unlock()
	if l.inited {
		return nil
	}
	l.Opts = &Options{
		LogFileDir:    "",
		AppName:       "app_log",
		ErrorFileName: "error.log",
		WarnFileName:  "warn.log",
		InfoFileName:  "info.log",
		DebugFileName: "debug.log",
		Level:         zapcore.DebugLevel,
		MaxSize:       100,
		MaxBackups:    60,
		MaxAge:        30,
	}
	if l.Opts.LogFileDir == "" {
		l.Opts.LogFileDir, _ = filepath.Abs(filepath.Dir(filepath.Join(".")))
		l.Opts.LogFileDir += sp + "logs" + sp
	}
	if l.Opts.Development {
		l.zapConfig = zap.NewDevelopmentConfig()
		l.zapConfig.EncoderConfig.EncodeTime = timeEncoder
	} else {
		l.zapConfig = zap.NewProductionConfig()
		l.zapConfig.EncoderConfig.EncodeTime = timeUnixNano
	}
	if l.Opts.OutputPaths == nil || len(l.Opts.OutputPaths) == 0 {
		l.zapConfig.OutputPaths = []string{"stdout"}
	}
	if l.Opts.ErrorOutputPaths == nil || len(l.Opts.ErrorOutputPaths) == 0 {
		l.zapConfig.OutputPaths = []string{"stderr"}
	}
	for _, fn := range mod {
		fn(l.Opts)
	}
	l.zapConfig.Level.SetLevel(l.Opts.Level)
	l.init()
	l.inited = true
	return l.Logger
}

func (l *Logger) init() {
	l.setSyncers()
	var err error
	l.Logger, err = l.zapConfig.Build(l.cores())
	if err != nil {
		panic(err)
	}
	defer func(Logger *zap.Logger) {
		err := Logger.Sync()
		if err != nil {
		}
	}(l.Logger)
}

func (l *Logger) setSyncers() {
	f := func(fN string) zapcore.WriteSyncer {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   l.Opts.LogFileDir + sp + l.Opts.AppName + "-" + fN,
			MaxSize:    l.Opts.MaxSize,
			MaxBackups: l.Opts.MaxBackups,
			MaxAge:     l.Opts.MaxAge,
			Compress:   true,
			LocalTime:  true,
		})
	}
	errWS = f(l.Opts.ErrorFileName)
	warnWS = f(l.Opts.WarnFileName)
	infoWS = f(l.Opts.InfoFileName)
	debugWS = f(l.Opts.DebugFileName)
	return
}

// SetMaxSize 最多文件个数
func SetMaxSize(MaxSize int) ModOptions {
	return func(option *Options) {
		option.MaxSize = MaxSize
	}
}

// SetMaxBackups 最多存在多少个切片文件
func SetMaxBackups(MaxBackups int) ModOptions {
	return func(option *Options) {
		option.MaxBackups = MaxBackups
	}
}

// SetMaxAge 保存的最大天数
func SetMaxAge(MaxAge int) ModOptions {
	return func(option *Options) {
		option.MaxAge = MaxAge
	}
}

// SetLogFileDir 设置日志目录
func SetLogFileDir(LogFileDir string) ModOptions {
	return func(option *Options) {
		option.LogFileDir = LogFileDir
	}
}

// SetAppName 设置应用名称
func SetAppName(AppName string) ModOptions {
	return func(option *Options) {
		option.AppName = AppName
	}
}

// SetLevel 设置日志级别
func SetLevel(Level zapcore.Level) ModOptions {
	return func(option *Options) {
		option.Level = Level
	}
}

// SetErrorFileName 设置Error级别日志文件名
func SetErrorFileName(ErrorFileName string) ModOptions {
	return func(option *Options) {
		option.ErrorFileName = ErrorFileName
	}
}

// SetWarnFileName 设置Warn级别日志文件名
func SetWarnFileName(WarnFileName string) ModOptions {
	return func(option *Options) {
		option.WarnFileName = WarnFileName
	}
}

// SetInfoFileName 设置Info级别日志文件名
func SetInfoFileName(InfoFileName string) ModOptions {
	return func(option *Options) {
		option.InfoFileName = InfoFileName
	}
}

// SetDebugFileName 设置Debug级别日志文件名
func SetDebugFileName(DebugFileName string) ModOptions {
	return func(option *Options) {
		option.DebugFileName = DebugFileName
	}
}

// SetDevelopment 设置开发者模式
func SetDevelopment(Development bool) ModOptions {
	return func(option *Options) {
		option.Development = Development
	}
}

func (l *Logger) cores() zap.Option {
	fileEncoder := zapcore.NewJSONEncoder(l.zapConfig.EncoderConfig)
	//consoleEncoder := zapcore.NewConsoleEncoder(l.zapConfig.EncoderConfig)
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	errPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.ErrorLevel && zapcore.ErrorLevel-l.zapConfig.Level.Level() > -1
	})
	warnPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.WarnLevel && zapcore.WarnLevel-l.zapConfig.Level.Level() > -1
	})
	infoPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.InfoLevel && zapcore.InfoLevel-l.zapConfig.Level.Level() > -1
	})
	debugPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl == zapcore.DebugLevel && zapcore.DebugLevel-l.zapConfig.Level.Level() > -1
	})
	cores := []zapcore.Core{
		zapcore.NewCore(fileEncoder, errWS, errPriority),
		zapcore.NewCore(fileEncoder, warnWS, warnPriority),
		zapcore.NewCore(fileEncoder, infoWS, infoPriority),
		zapcore.NewCore(fileEncoder, debugWS, debugPriority),
	}
	if l.Opts.Development {
		cores = append(cores, []zapcore.Core{
			zapcore.NewCore(consoleEncoder, errorConsoleWS, errPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, warnPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, infoPriority),
			zapcore.NewCore(consoleEncoder, debugConsoleWS, debugPriority),
		}...)
	}
	return zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return zapcore.NewTee(cores...)
	})
}
func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}

func timeUnixNano(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendInt64(t.UnixNano() / 1e6)
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, fields ...zapcore.Field) {
	l.Debug(msg, fields...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func Debugf(template string, args ...interface{}) {
	zLogger.Debugf(template, args)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
//
// When debug-level logging is disabled, this is much faster than
//  s.With(keysAndValues).Debug(msg)
func Debugw(msg string, keysAndValues ...interface{}) {
	zLogger.Debugw(msg, keysAndValues)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, fields ...zapcore.Field) {
	l.Info(msg, fields...)
}

// Infof uses fmt.Sprintf to log a templated message.
func Infof(template string, args ...interface{}) {
	zLogger.Infof(template, args...)
}

// Infow logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Infow(msg string, keysAndValues ...interface{}) {
	zLogger.Infow(msg, keysAndValues)
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Warn(msg string, fields ...zapcore.Field) {
	l.Warn(msg, fields...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func Warnf(template string, args ...interface{}) {
	zLogger.Warnf(template, args...)
}

// Warnw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Warnw(msg string, keysAndValues ...interface{}) {
	zLogger.Warnw(msg, keysAndValues)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, fields ...zapcore.Field) {
	l.Error(msg, fields...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func Errorf(template string, args ...interface{}) {
	zLogger.Errorf(template, args...)
}

// Errorw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in With.
func Errorw(msg string, keysAndValues ...interface{}) {
	zLogger.Errorw(msg, keysAndValues)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, fields ...zapcore.Field) {
	l.Panic(msg, fields...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func Panicf(template string, args ...interface{}) {
	zLogger.Panicf(template, args...)
}

// Panicw logs a message with some additional context, then panics. The
// variadic key-value pairs are treated as they are in With.
func Panicw(msg string, keysAndValues ...interface{}) {
	zLogger.Panicw(msg, keysAndValues)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, fields ...zapcore.Field) {
	l.Fatal(msg, fields...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func Fatalf(template string, args ...interface{}) {
	zLogger.Fatalf(template, args...)
}

// Fatalw logs a message with some additional context, then calls os.Exit. The
// variadic key-value pairs are treated as they are in With.
func Fatalw(msg string, keysAndValues ...interface{}) {
	zLogger.Fatalw(msg, keysAndValues)
}
