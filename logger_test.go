package log

import (
	"fmt"
	"testing"
)

func Test_LoggerCountsErrorsAndWarnings(t *testing.T) {
	ConfigureMainLoggerForTest(t)
	for i := 0; i < 10; i++ {
		Error(fmt.Sprintf("%d", i))
		Warning(fmt.Sprintf("%d", i))
	}
	if numberOfErrorsLogged != 10 {
		t.Errorf(
			"expected numberOfErrorsLogged to be 10 but was %d",
			numberOfErrorsLogged,
		)
	}
	if numberOfWarningsLogged != 10 {
		t.Errorf(
			"expected numberOfWarningsLogged to be 10 but was %d",
			numberOfWarningsLogged,
		)
	}
}

func Test_ErrorAmountReturnsCorrectNumber(t *testing.T) {
	ConfigureMainLoggerForTest(t)
	Info("")
	Warning("")
	Error("")
	if ErrorAmount() != numberOfErrorsLogged {
		t.Errorf(
			"expected ErrorAmoount to return %d but returned %d",
			numberOfErrorsLogged,
			ErrorAmount(),
		)
	}
}

func Test_WarningAmountReturnsCorrectNumber(t *testing.T) {
	ConfigureMainLoggerForTest(t)
	Info("")
	Warning("")
	Error("")
	if WarningAmount() != numberOfWarningsLogged {
		t.Errorf(
			"expected ErrorAmoount to return %d but returned %d",
			numberOfWarningsLogged,
			WarningAmount(),
		)
	}
}
