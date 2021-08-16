package log

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

const directoryDateFormat = "2006_01_02"

var (
	log                    *zap.Logger
	detailedLog            *zap.Logger
	compilationLog         *zap.Logger
	numberOfErrorsLogged   = 0
	numberOfWarningsLogged = 0
)

// ConfigureMainLoggerForTest configures the main logger to be used during a test.
func ConfigureMainLoggerForTest(t *testing.T) {
	log = zaptest.NewLogger(t, zaptest.Level(zap.WarnLevel))
	compilationLog = zaptest.NewLogger(t, zaptest.Level(zap.WarnLevel))
}

// ConfigureMainLogger configures a logger and creates a log file in the specified
// directory, but under a second directory named after the current date.
func ConfigureMainLogger(directory string, fileName string) error {
	if err := configureCompilationLogger(directory, fileName); err != nil {
		return err
	}
	encoding := zap.NewDevelopmentConfig().Encoding
	config := zap.Config{
		OutputPaths: []string{
			createLogPath(
				directory,
				fmt.Sprintf("%s.log", fileName),
			),
		},
		Encoding: encoding,
		Level:    zap.NewAtomicLevelAt(zap.InfoLevel),
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "msg",
			LevelKey:     "level",
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	var err error
	log, err = config.Build()
	return err
}

// ConfigureDetailedLoggerForTest configures the detailed logger to be used during
// a test.
func ConfigureDetailedLoggerForTest(t *testing.T) {
	detailedLog = zaptest.NewLogger(t, zaptest.Level(zap.WarnLevel))
}

// ConfigureDetailedLogger creates a detailed logger that will log everything that
// the main logger logs as well as any logs created with Detailed* functions.
func ConfigureDetailedLogger(directory string, fileName string) error {
	if log == nil {
		return errors.New("must configure main logger before detailed logger")
	}
	encoding := zap.NewDevelopmentConfig().Encoding
	config := zap.Config{
		OutputPaths: []string{
			createLogPath(
				directory,
				fmt.Sprintf("%s_DETAILED.log", fileName),
			),
		},
		Encoding: encoding,
		Level:    zap.NewAtomicLevelAt(zap.DebugLevel),
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "msg",
			LevelKey:     "level",
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	var err error
	detailedLog, err = config.Build()
	return err
}

// configureCompilationLogger creates a logger that will log everything that the
// detailed and main logger logs in a format that can be interpreted by a future
// compilation system for logs.
func configureCompilationLogger(directory string, fileName string) error {
	if !logDirectoryExists(directory) {
		if err := createLogDirectory(directory); err != nil {
			return err
		}
	}
	encoding := zap.NewProductionConfig().Encoding
	config := zap.Config{
		OutputPaths: []string{
			createLogPath(
				directory,
				fmt.Sprintf("%s_COMPILATION.log", fileName),
			),
		},
		Encoding: encoding,
		Level:    zap.NewAtomicLevelAt(zap.DebugLevel),
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "msg",
			LevelKey:     "level",
			TimeKey:      "time",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
	var err error
	compilationLog, err = config.Build()
	return err
}

// logDirectoryExists returns true if the standard log directory defined in the
// configuration file, and a sub-directory named after the current days date, exists.
func logDirectoryExists(directory string) bool {
	_, err := os.Stat(todaysLogDirectoryPath(directory))
	return !os.IsNotExist(err)
}

// createLogDirectory creates each missing directory on the path to the current
// log directory.
func createLogDirectory(directory string) error {
	return os.MkdirAll(todaysLogDirectoryPath(directory), 0755)
}

// todaysLogDirectoryPath returns the path to all logs that have been taken during
// the current day
func todaysLogDirectoryPath(directory string) string {
	today := time.Now()
	return fmt.Sprintf("%s/%s", directory, today.Format(directoryDateFormat))
}

// createLogPath joins together the log directory path with the logs filename
func createLogPath(directory, fileName string) string {
	return fmt.Sprintf("%s/%s", todaysLogDirectoryPath(directory), fileName)
}

// Field returns a zap.Field which can be used for adding custom properties to
// a log entry.
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// ErrorField returns a zap.Field which is pre-configured with the key "cause" and
// the value being the error message from `err`
func ErrorField(err error) zap.Field {
	return Field("cause", err.Error())
}

// ObjectField is meant to take a struct type and marshal it into a zap.Field, which
// it returns
func ObjectField(key string, value interface{}) zap.Field {
	marshaler := zapcore.ObjectMarshalerFunc(func(oe zapcore.ObjectEncoder) error {
		buildObjectField(oe, value)
		return nil
	})
	return zap.Object(key, marshaler)
}

// buildObjectField goes through the fields of an arbitrary interface and adds the
// field name and value unto the zap encoder
func buildObjectField(enc zapcore.ObjectEncoder, value interface{}) {
	v := reflect.Indirect(reflect.ValueOf(value))
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		fn := v.Type().Field(i).Name
		switch f.Interface().(type) {
		case string:
			enc.AddString(fn, f.String())
		case int, int64:
			enc.AddInt64(fn, f.Int())
		case float64:
			enc.AddFloat64(fn, f.Float())
		case bool:
			enc.AddBool(fn, f.Bool())
		case *time.Time:
			enc.AddString(fn, formatDate(f.Interface().(*time.Time)))
		}
	}
}

// formatDate returns a string containing the default log representation
// of a date value
func formatDate(date *time.Time) string {
	format := "2006-01-02"
	if date == nil {
		return "nil"
	} else {
		return date.Format(format)
	}
}

// Info places an info log entry in both the general and the detailed (debug) log.
func Info(msg string, fields ...zap.Field) {
	log.Info(fmt.Sprintf(" %s", msg), fields...)
	compilationLog.Info(msg, fields...)
	if detailedLog != nil {
		detailedLog.Info(fmt.Sprintf(" %s", msg), fields...)
	}
}

// DetailInfo places an info log entry in only the detailed (debug) log
func DetailInfo(msg string, fields ...zap.Field) {
	detailedLog.Info(fmt.Sprintf(" %s", msg), fields...)
	compilationLog.Info(msg, fields...)
}

// Error places an error log entry in both the general and the detailed (debug) log
// as well as increments the error counter
func Error(msg string, fields ...zap.Field) {
	numberOfErrorsLogged++
	if _, filePath, line, ok := runtime.Caller(1); ok {
		splitPath := strings.Split(filePath, "/")
		file := strings.Join(splitPath[len(splitPath)-2:], "/")
		log.Error(fmt.Sprintf(" ocurred at %s, line %d: %s", file, line, msg), fields...)
		compilationLog.Error(fmt.Sprintf("ocurred at %s, line %d: %s", file, line, msg), fields...)
		if detailedLog != nil {
			detailedLog.Error(fmt.Sprintf(" ocurred at %s, line %d: %s", file, line, msg), fields...)
		}
	} else {
		log.Error(fmt.Sprintf(" %s", msg), fields...)
		compilationLog.Error(fmt.Sprintf(" %s", msg), fields...)
		if detailedLog != nil {
			detailedLog.Error(fmt.Sprintf(" %s", msg), fields...)
		}
	}
}

// Warning places an warning log entry in both the general and the detailed (debug) log
// as well as increments the warnings counter
func Warning(msg string, fields ...zap.Field) {
	numberOfWarningsLogged++
	log.Warn(fmt.Sprintf(" %s", msg), fields...)
	compilationLog.Warn(msg, fields...)
	if detailedLog != nil {
		detailedLog.Warn(fmt.Sprintf(" %s", msg), fields...)
	}
}

// Flush checks if a detailed log has been configured, if it has then it is flushed to
// the disk and a potential error is returned to the caller. The main log is also
// flushed and any potential error resulting from that is also returned to the caller.
func Flush() error {
	if detailedLog != nil {
		if err := detailedLog.Sync(); err != nil {
			return err
		}
	}
	if err := compilationLog.Sync(); err != nil {
		return err
	}
	return log.Sync()
}

// ErrorAmount returns the number of errors that has been logged at the time of
// calling NumberOfErrors. This method is used as opposed to exporting the
// numberOfErrorsLogged variable to make sure that the variable cannot be mutated
// outside of the logger package.
func ErrorAmount() int {
	return numberOfErrorsLogged
}

// NumberOfWarnings returns the number of warnings that has been logged at the time
// of calling NumberOfErrors. This method is used as opposed to exporting the
// numberOfWarningsLogged variable to make sure that the variable cannot be mutated
// outside of the logger package.
func WarningAmount() int {
	return numberOfWarningsLogged
}
