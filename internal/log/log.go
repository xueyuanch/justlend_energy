package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
)

var (
	// Enable all logger levels.
	allLevels = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.DebugLevel
	})
)

// fileWrite constructing a lumberjack.Logger to creates & writes logs to a file
// and will help us to roll log files automatically.
func fileWriter(fName string) io.Writer {
	// Customize io writer and use lumberjack to rolling log files.
	// For more about lumberjack, check: https://github.com/natefinch/lumberjack
	return &lumberjack.Logger{
		MaxSize:    200, // Max size in megabytes.
		MaxBackups: 30,  // Max number of backup files for retention.
		MaxAge:     30,  // Retention age in days.
	}
}

// NewLogger is a factory to construct a logger with the given name. If the
// application runs in release mode, it will create a log file with the passed
// name in the configured log directory and write all subsequent logs of the
// logger to the file. Otherwise, the logs will be output to the standard output console.
//
// Note: The conflict of loggers with the same name has not been handled, consumer
// should avoid using duplicate names.
func NewLogger(name string) *zap.Logger {
	// Directly print the log to the standard output console in the local
	// environment for the convenience of debugging
	return zap.New(zapcore.NewCore(zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.Lock(os.Stdout), allLevels,
	))
}

// Initialize a default package level sugared logger.
var sugar = NewLogger("energy").Sugar()

func Debug(args ...interface{})                           { sugar.Debug(args...) }
func DebugF(template string, args ...interface{})         { sugar.Debugf(template, args...) }
func DebugW(message string, keysAndValues ...interface{}) { sugar.Debugw(message, keysAndValues...) }

func Info(args ...interface{})                           { sugar.Info(args...) }
func InfoF(template string, args ...interface{})         { sugar.Infof(template, args...) }
func InfoW(message string, keysAndValues ...interface{}) { sugar.Infow(message, keysAndValues...) }

func Warn(args ...interface{})                           { sugar.Warn(args...) }
func WarnF(template string, args ...interface{})         { sugar.Warnf(template, args...) }
func WarnW(message string, keysAndValues ...interface{}) { sugar.Warnw(message, keysAndValues...) }

func Error(args ...interface{})                           { sugar.Error(args...) }
func Errorf(template string, args ...interface{})         { sugar.Errorf(template, args...) }
func ErrorW(message string, keysAndValues ...interface{}) { sugar.Errorw(message, keysAndValues...) }

func DPanic(args ...interface{})                           { sugar.DPanic(args...) }
func DPanicF(template string, args ...interface{})         { sugar.DPanicf(template, args...) }
func DPanicW(message string, keysAndValues ...interface{}) { sugar.DPanicw(message, keysAndValues...) }

func Panic(args ...interface{})                           { sugar.Panic(args...) }
func PanicF(template string, args ...interface{})         { sugar.Panicf(template, args...) }
func PanicW(message string, keysAndValues ...interface{}) { sugar.Panicw(message, keysAndValues...) }

func Fatal(args ...interface{})                           { sugar.Fatal(args...) }
func Fatalf(template string, args ...interface{})         { sugar.Fatalf(template, args...) }
func FatalW(message string, keysAndValues ...interface{}) { sugar.Fatalw(message, keysAndValues...) }
