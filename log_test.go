package log

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func Test_LogAttachesChildren(t *testing.T) {
	parent := &log{
		logger:   zaptest.NewLogger(t, zaptest.Level(zap.FatalLevel)),
		attached: make([]Log, 0),
	}
	child := &LogMock{
		InfoMock: func(s string, f ...zap.Field) {},
	}
	parent.Attach(child)
	if len(parent.attached) != 1 {
		t.Errorf(
			"expected number of attached to be 1 but found %d",
			len(parent.attached),
		)
	}
}

func Test_InfoForwardsToChildren(t *testing.T) {
	parent := &log{
		logger:   zaptest.NewLogger(t, zaptest.Level(zap.FatalLevel)),
		attached: make([]Log, 0),
	}
	child := &LogMock{
		InfoMock: func(s string, f ...zap.Field) {},
	}
	parent.Attach(child)
	parent.Info("")
	if child.InfoCalls != 1 {
		t.Errorf(
			"expected 1 call to childs Info method, found %d",
			child.InfoCalls,
		)
	}
}

func Test_ErrorForwardsToChildren(t *testing.T) {
	parent := &log{
		logger:   zaptest.NewLogger(t, zaptest.Level(zap.FatalLevel)),
		attached: make([]Log, 0),
	}
	child := &LogMock{
		ErrorMock: func(s string, f ...zap.Field) {},
	}
	parent.Attach(child)
	parent.Error("")
	if child.ErrorCalls != 1 {
		t.Errorf(
			"expected 1 call to childs Error method, found %d",
			child.ErrorCalls,
		)
	}
}

func Test_WarnForwardsToChildren(t *testing.T) {
	parent := &log{
		logger:   zaptest.NewLogger(t, zaptest.Level(zap.FatalLevel)),
		attached: make([]Log, 0),
	}
	child := &LogMock{
		WarnMock: func(s string, f ...zap.Field) {},
	}
	parent.Attach(child)
	parent.Warn("")
	if child.WarnCalls != 1 {
		t.Errorf(
			"expected 1 call to childs Warn method, found %d",
			child.WarnCalls,
		)
	}
}

func Test_FlushForwardsToChildren(t *testing.T) {
	parent := &log{
		logger:   zaptest.NewLogger(t, zaptest.Level(zap.FatalLevel)),
		attached: make([]Log, 0),
	}
	child := &LogMock{
		FlushMock: func() error {
			return nil
		},
	}
	parent.Attach(child)
	parent.Flush()
	if child.FlushCalls != 1 {
		t.Errorf(
			"expected 1 call to childs Flush method, found %d",
			child.FlushCalls,
		)
	}
}

func Test_CountsErrors(t *testing.T) {
	logger := &log{
		logger:   zaptest.NewLogger(t, zaptest.Level(zap.FatalLevel)),
		attached: make([]Log, 0),
	}
	logger.Error("1")
	if logger.errors != 1 {
		t.Errorf(
			"expected errors to be 1 but was %d",
			logger.errors,
		)
	}
}

func Test_CountsWarnings(t *testing.T) {
	logger := &log{
		logger:   zaptest.NewLogger(t, zaptest.Level(zap.FatalLevel)),
		attached: make([]Log, 0),
	}
	logger.Warn("1")
	if logger.warnings != 1 {
		t.Errorf(
			"expected warnings to be 1 but was %d",
			logger.warnings,
		)
	}
}

func Test_ErrorsReturnsCorrectAmount(t *testing.T) {
	logger := &log{
		logger:   zaptest.NewLogger(t, zaptest.Level(zap.FatalLevel)),
		attached: make([]Log, 0),
	}
	logger.Error("1")
	if logger.Errors() != 1 {
		t.Errorf(
			"expected Errors to return 1 but was %d",
			logger.Errors(),
		)
	}
}

func Test_WarningsReturnsCorrectAmount(t *testing.T) {
	logger := &log{
		logger:   zaptest.NewLogger(t, zaptest.Level(zap.FatalLevel)),
		attached: make([]Log, 0),
	}
	logger.Warn("1")
	if logger.Warnings() != 1 {
		t.Errorf(
			"expected warnings to be 1 but was %d",
			logger.Warnings(),
		)
	}
}

func Test_PathReturnsCorrectString(t *testing.T) {
	logger := &log{
		directory: "/usr/share/logs",
		fileName:  "test.log",
	}
	if logger.Path() != "/usr/share/logs/test.log" {
		t.Errorf(
			"expected Path to return '/usr/share/logs/test.log' but found '%s'",
			logger.Path(),
		)
	}
}
