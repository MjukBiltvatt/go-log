package log

import "go.uber.org/zap"

type LogMock struct {
	InfoMock  func(string, ...zap.Field)
	InfoCalls int

	ErrorMock  func(string, ...zap.Field)
	ErrorCalls int

	ErrorsMock  func() int
	ErrorsCalls int

	WarnMock  func(string, ...zap.Field)
	WarnCalls int

	WarningsMock  func() int
	WarningsCalls int

	AttachMock  func(Log)
	AttachCalls int

	FlushMock  func() error
	FlushCalls int

	PathMock  func() string
	PathCalls int
}

func (mock *LogMock) Info(msg string, fields ...zap.Field) {
	mock.InfoCalls++
	mock.InfoMock(msg, fields...)
}

func (mock *LogMock) Error(msg string, fields ...zap.Field) {
	mock.ErrorCalls++
	mock.ErrorMock(msg, fields...)
}

func (mock *LogMock) Errors() int {
	mock.ErrorsCalls++
	return mock.ErrorsMock()
}

func (mock *LogMock) Warn(msg string, fields ...zap.Field) {
	mock.WarnCalls++
	mock.WarnMock(msg, fields...)
}

func (mock *LogMock) Warnings() int {
	mock.WarningsCalls++
	return mock.WarningsMock()
}

func (mock *LogMock) Attach(log Log) {
	mock.AttachCalls++
	mock.AttachMock(log)
}

func (mock *LogMock) Flush() error {
	mock.FlushCalls++
	return mock.FlushMock()
}

func (mock *LogMock) Path() string {
	mock.PathCalls++
	return mock.PathMock()
}
