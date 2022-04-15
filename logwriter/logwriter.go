package logwriter

import (
	"locker/constants"
	"os"
)

type LogWriter interface {
	Write(msg string, level int)
	GetLevel() int
}

type FileWriter struct {
	level int
	path  string
}

type TerminalWriter struct {
	level int
}

type CombinedWriter struct {
	lowestRelevantLevel int
	loggerList          []LogWriter
}

func DefaultDebugWriter() LogWriter {
	return NewTerminalWriter(constants.LevelDebug)
}

func DefaultProdWriter() LogWriter {
	stdWriter := NewTerminalWriter(constants.LevelInfo)
	fileWriter := NewFileWriter(constants.LevelError, "./errors.log")
	list := make([]LogWriter, 2)
	list[0] = stdWriter
	list[1] = fileWriter
	return NewCombinedWriter(list)
}

func NewCombinedWriter(loggers []LogWriter) CombinedWriter {
	lowest := 9999
	for _, v := range loggers {
		if v.GetLevel() < lowest {
			lowest = v.GetLevel()
		}
	}
	return CombinedWriter{
		lowestRelevantLevel: lowest,
		loggerList:          loggers,
	}
}

func (c CombinedWriter) Write(msg string, level int) {
	for _, v := range c.loggerList {
		v.Write(msg, level)
	}
}

func (c CombinedWriter) GetLevel() int {
	return c.lowestRelevantLevel
}

func NewTerminalWriter(level int) TerminalWriter {
	return TerminalWriter{}
}

func (t TerminalWriter) Write(msg string, level int) {
	if t.level <= level {
		if t.level == constants.LevelError {
			_, _ = os.Stderr.WriteString(msg)
		}
		_, _ = os.Stdout.WriteString(msg)
	}
}

func (t TerminalWriter) GetLevel() int {
	return t.level
}

func NewFileWriter(level int, path string) FileWriter {
	return FileWriter{
		level: level,
		path:  path,
	}
}

func (f FileWriter) Write(msg string, level int) {
	if f.level <= level {
		file, err := os.OpenFile(f.path, os.O_APPEND, 0644)
		if err != nil {
			panic(any(err))
		}

		file.WriteString(msg)
		file.Close()
	}
}

func (f FileWriter) GetLevel() int {
	return f.level
}
