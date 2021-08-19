package log

import "go.uber.org/zap"

type logMock struct {
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

func (mock *logMock) Info(msg string, fields ...zap.Field) {
	mock.InfoCalls++
	mock.InfoMock(msg, fields...)
}

func (mock *logMock) Error(msg string, fields ...zap.Field) {
	mock.ErrorCalls++
	mock.ErrorMock(msg, fields...)
}

func (mock *logMock) Errors() int {
	mock.ErrorsCalls++
	return mock.ErrorsMock()
}

func (mock *logMock) Warn(msg string, fields ...zap.Field) {
	mock.WarnCalls++
	mock.WarnMock(msg, fields...)
}

func (mock *logMock) Warnings() int {
	mock.WarningsCalls++
	return mock.WarningsMock()
}

func (mock *logMock) Attach(log Log) {
	mock.AttachCalls++
	mock.AttachMock(log)
}

func (mock *logMock) Flush() error {
	mock.FlushCalls++
	return mock.FlushMock()
}

func (mock *logMock) Path() string {
	mock.PathCalls++
	return mock.PathMock()
}
