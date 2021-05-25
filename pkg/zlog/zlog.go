package zlog

import (
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"time"
)

var zLogger *zap.SugaredLogger

func init() {

	createDirectoryIfNotExists()

	// 实现两个判断日志等级的interface
	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.InfoLevel
	})

	debugLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})

	errorLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})

	// 获取info, error 日志文件的io.Writer
	infoWriter := getWriter("k8s_aim_info.log")
	debugWriter := getWriter("k8s_aim_debug.log")
	errorWriter := getWriter("kus_aim_error.log")

	encoder := getEncoder()
	core := zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWriter), errorLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(debugWriter), debugLevel),
	)

	log := zap.New(core, zap.AddCaller())
	zLogger = log.Sugar()

}

func createDirectoryIfNotExists() {
	path, _ := os.Getwd()

	if _, err := os.Stat(fmt.Sprintf("%s/logs", path)); os.IsNotExist(err) {
		_ = os.Mkdir("logs", os.ModePerm)
	}
}

func getWriter(fileName string) io.Writer {
	path, _ := os.Getwd()
	fileName = fmt.Sprintf("%s/logs/%s", path, fileName)
	jackLogger := &lumberjack.Logger{
		Filename:   strings.Replace(fileName, ".log", "", -1) + "-%Y%m%d%H.log",
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return zapcore.AddSync(jackLogger)
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05"))
	}
	encoderConfig.TimeKey = "ts"
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeDuration = func(duration time.Duration, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendInt64(int64(duration) / 1000000)
	}
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func Debugf(template string, args ...interface{}) {
	zLogger.Debugf(template, args)
}

func Infof(template string, args ...interface{}) {
	zLogger.Infof(template, args...)
}

func Warn(args ...interface{}) {
	zLogger.Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	zLogger.Warnf(template, args...)
}

func Error(args ...interface{}) {
	zLogger.Error(args...)
}

func Errorf(template string, args ...interface{}) {
	zLogger.Errorf(template, args...)
}

func DPanic(args ...interface{}) {
	zLogger.DPanic(args...)
}

func DPanicf(template string, args ...interface{}) {
	zLogger.DPanicf(template, args...)
}

func Panic(args ...interface{}) {
	zLogger.Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	zLogger.Panicf(template, args...)
}

func Fatal(args ...interface{}) {
	zLogger.Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	zLogger.Fatalf(template, args...)
}
