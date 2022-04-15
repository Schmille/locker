package locker

import (
	"bytes"
	"locker/constants"
	"locker/logwriter"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

type context struct {
	lock     sync.RWMutex
	contexts map[string]*contextValues
}

type contextValues struct {
	context   []string
	generated string
}

var writer = logwriter.DefaultDebugWriter()
var goroutineContexts = context{
	lock:     sync.RWMutex{},
	contexts: make(map[string]*contextValues),
}

func SetWriter(logwriter logwriter.LogWriter) {
	writer = logwriter
}

func Debug(msg string) {
	writeWithLevel(msg, constants.LevelDebug)
}

func Info(msg string) {
	writeWithLevel(msg, constants.LevelInfo)
}

func Warn(msg string) {
	writeWithLevel(msg, constants.LevelWarn)
}

func Error(msg string) {
	writeWithLevel(msg, constants.LevelError)
}

func Push(context string) {
	goroutineContexts.ensureValuesExists()
	goroutineContexts.lock.Lock()
	goroutineContexts.contexts[getGoRoutineID()].push(context)
	goroutineContexts.lock.Unlock()
}

func Pop() {
	goroutineContexts.lock.Lock()
	goroutineContexts.contexts[getGoRoutineID()].pop()
	goroutineContexts.lock.Unlock()
}

func writeWithLevel(msg string, level int) {
	sb := strings.Builder{}
	sb.WriteString(getFormattedTimestamp())
	sb.WriteString(" ")
	sb.WriteString(getGoRoutineID())
	sb.WriteString(" ")
	sb.WriteString(getLevelFormat(level))
	sb.WriteString(" ")
	sb.WriteString(goroutineContexts.GetContext())
	sb.WriteString(msg)
	sb.WriteString("\n")
	writer.Write(sb.String(), level)
}

func (c context) GetContext() string {
	if !c.valuesExists() {
		return ""
	}

	id := getGoRoutineID()
	c.lock.RLock()
	result := c.contexts[id].get()
	c.lock.RUnlock()
	return result
}

func (c context) ensureValuesExists() {
	id := getGoRoutineID()
	if !c.valuesExists() {
		c.lock.Lock()
		c.contexts[id] = &contextValues{
			context:   make([]string, 0),
			generated: "",
		}

		c.lock.Unlock()
	}
}

func (c context) valuesExists() bool {
	_, present := c.contexts[getGoRoutineID()]
	return present
}

func (c *contextValues) push(ctx string) {
	c.context = append(c.context, ctx)
	c.regenerate()
}

func (c *contextValues) pop() {
	if len(c.context) == 0 {
		return
	}
	c.context = c.context[:len(c.context)-1]
	c.regenerate()
}

func (c *contextValues) regenerate() {
	builder := strings.Builder{}
	for _, v := range c.context {
		builder.WriteString("[")
		builder.WriteString(v)
		builder.WriteString("] ")
	}
	c.generated = builder.String()
}

func (c contextValues) get() string {
	return c.generated
}

func getLevelFormat(level int) string {
	switch level {
	case constants.LevelDebug:
		return "[DEBUG]"
		break
	case constants.LevelInfo:
		return "[INFO ]"
		break
	case constants.LevelWarn:
		return "[WARN ]"
		break
	case constants.LevelError:
		return "[ERROR]"
		break
	}
	return "[UNKWN]"
}

func getFormattedTimestamp() string {
	return time.Now().Format("[2006-01-02-15:04:05]")
}

func getGoRoutineID() string {
	gr := bytes.Fields(debug.Stack())[1]
	return "[routineno. " + string(gr) + "]"
}