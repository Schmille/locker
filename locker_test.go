package locker

import (
	"locker/constants"
	"strings"
	"testing"
)

type MockLogWriter struct {
	Msg      string
	variable *string
}

func (m MockLogWriter) Write(msg string, level int) {
	m.Msg = msg
	*m.variable = m.Msg
}

func (m MockLogWriter) GetLevel() int {
	return constants.LevelDebug
}

func NewMockWriter(variable *string) MockLogWriter {
	return MockLogWriter{
		Msg:      "",
		variable: variable,
	}
}

func TestDebug(t *testing.T) {
	var result string
	SetWriter(NewMockWriter(&result))
	Debug("Debug Test Message!")
	if !strings.HasSuffix(result, "[DEBUG] Debug Test Message!\n") {
		t.Fail()
	}
}

func TestError(t *testing.T) {
	var result string
	SetWriter(NewMockWriter(&result))
	Error("Error Test Message!")
	if !strings.HasSuffix(result, "[ERROR] Error Test Message!\n") {
		t.Fail()
	}
}

func TestInfo(t *testing.T) {
	var result string
	SetWriter(NewMockWriter(&result))
	Info("Info Test Message!")
	if !strings.HasSuffix(result, "[INFO ] Info Test Message!\n") {
		t.Fail()
	}
}

func TestPop(t *testing.T) {
	var result string
	SetWriter(NewMockWriter(&result))

	Push("ThisShouldNotBeHere")
	Pop()

	Debug("Test Message!")
	if strings.Contains(result, "ThisShouldNotBeHere") {
		t.Fail()
	}
}

func TestPush(t *testing.T) {
	var result string
	SetWriter(NewMockWriter(&result))
	Push("SomeContext")
	Debug("Debug Test Message!")
	if !strings.HasSuffix(result, "[DEBUG] [SomeContext] Debug Test Message!\n") {
		t.Fail()
	}
}

func TestWarn(t *testing.T) {
	var result string
	SetWriter(NewMockWriter(&result))
	Warn("Warn Test Message!")
	if !strings.HasSuffix(result, "[WARN ] Warn Test Message!\n") {
		t.Fail()
	}
}
