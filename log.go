package log

import (
	"fmt"

	"go.uber.org/zap"
)

// Log contains the methods required for each logger, it is left exported
// so that users of the package can implement mocks as needed.
type Log interface {
	// Info logs an info entry which consists of an informative message
	// and optional fields. To log such fields, please refer to the
	// Field() function
	Info(string, ...zap.Field)
	// Error logs an error entry which consists of an informative message
	// and optional fields. To log such fields, please refer to the
	// Field() function
	Error(string, ...zap.Field)
	// Errors returns the amount of errors that has been logged at the
	// point of calling.
	Errors() int
	// Warn logs a warning entry which consists of an informative message
	// and optional fields. To log such fields, please refer to the
	// Field() function
	Warn(string, ...zap.Field)
	// Warnings returns the amount of warnings that has been logged at
	// the point of calling.
	Warnings() int
	// Attach registers a logger as a child logger to the receiving value.
	// This will cause any logs performed on the parent to be forwarded
	// to the child.
	Attach(Log)
	// Flush synchronizes the logger and its children to the harddrive and
	// thus forces a log flush to disk. Note that the request is forwarded
	// to any child loggers attached to the logger.
	Flush() error
	// Path returns the file path to the log in which entries are written.
	Path() string
}

// log is a concrete implementation of the Log interface which uses an
// underlying zap.Logger to perform its logging.
type log struct {
	logger    *zap.Logger
	attached  []Log
	errors    int
	warnings  int
	directory string
	fileName  string
}

func (log *log) Info(msg string, fields ...zap.Field) {
	log.logger.Info(msg, fields...)
	for _, sub := range log.attached {
		sub.Info(msg, fields...)
	}
}

func (log *log) Error(msg string, fields ...zap.Field) {
	log.logger.Error(msg, fields...)
	log.errors++
	for _, sub := range log.attached {
		sub.Error(msg, fields...)
	}
}

func (log *log) Errors() int {
	return log.errors
}

func (log *log) Warn(msg string, fields ...zap.Field) {
	log.logger.Warn(msg, fields...)
	log.warnings++
	for _, sub := range log.attached {
		sub.Warn(msg, fields...)
	}
}

func (log *log) Warnings() int {
	return log.warnings
}

func (log *log) Attach(sub Log) {
	log.attached = append(log.attached, sub)
}

func (log *log) Flush() error {
	if err := log.logger.Sync(); err != nil {
		return err
	}
	for _, sub := range log.attached {
		if err := sub.Flush(); err != nil {
			return err
		}
	}
	return nil
}

func (log *log) Path() string {
	return fmt.Sprintf("%s/%s", log.directory, log.fileName)
}
