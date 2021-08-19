package log

import (
	"fmt"
	"os"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
)

// NewTextLogger builds a plain-text based logger and returns it to the caller. If the
// log directory has not been created then this is created here. An error may be returned
// if there was a problem in building the underlying logger or if the directory could not
// be created.
func NewTextLogger(directory, name string) (Log, error) {
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(directory, 0755); err != nil {
			return nil, err
		}
	}
	config := zapConfig(
		"console",
		fmt.Sprintf("%s/%s", directory, name),
	)
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &log{
		logger:   logger,
		attached: make([]Log, 0),
	}, nil
}

// NewTestLogger returns a logger that has been configured for testing. None of the
// logs will be written anywhere and instead just discarded.
func NewTestLogger(t *testing.T) Log {
	return &log{
		logger:   zaptest.NewLogger(t, zaptest.Level(zap.FatalLevel)),
		attached: make([]Log, 0),
	}
}

// NewJSONLogger builds a JSON based logger and returns it to the caller. If the log
// directory has not been created then this is created here. An error may be returned
// if there was a problem in building the underlying logger or if the directory could
// not be created.
func NewJSONLogger(directory, name string) (Log, error) {
	_, err := os.Stat(directory)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(directory, 0755); err != nil {
			return nil, err
		}
	}
	config := zapConfig(
		"json",
		fmt.Sprintf("%s/%s", directory, name),
	)
	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	return &log{
		logger:   logger,
		attached: make([]Log, 0),
	}, nil
}

// zapConfig returns a default zap.Config that can be shared between loggers for a
// more consistent experience when reading them. The type of encoding and the
// output path must be provided as parameters
func zapConfig(encoding, filePath string) zap.Config {
	zap.NewProduction()
	return zap.Config{
		OutputPaths: []string{filePath},
		Encoding:    encoding,
		Level:       zap.NewAtomicLevelAt(zap.DebugLevel),
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "message",
			LevelKey:     "level",
			TimeKey:      "timestamp",
			EncodeTime:   zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
			EncodeLevel:  zapcore.CapitalLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}
}

// Field returns a zap.Field which can be used for adding custom properties to
// a log entry.
func Field(key string, value interface{}) zap.Field {
	return zap.Any(key, value)
}

// ErrField returns a zap.Field which is pre-configured with the key "cause" and
// with the value being the error message from the provided error.
func ErrField(err error) zap.Field {
	return Field("cause", err.Error())
}
