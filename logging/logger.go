package logging

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const (
	Version   = "1"
	Marker    = "*"
	Separator = "|"

	NOP uint8 = iota
	ALARM
	ERROR
	STAT
	INFO
	EVENT
	DEBUG
	UNKNOWN
)

// TODO build this out and make sure our default implementation follows it.
type Logger interface {
}

type severityWriter struct {
	writer io.Writer
	name   string
}

type DefaultLogger struct {
	name            string
	application     string
	group           string
	pid             string
	correlationId   string
	taskId          string
	severity        uint8
	severityStreams map[uint8]severityWriter
}

func NewDefaultLogger() *DefaultLogger {
	path := strings.Split(os.Args[0], "/")
	name := path[len(path)-1]
	application := os.Getenv("MON_APP")
	if application == "" {
		application = name
	}
	group := os.Getenv("MON_GROUP")
	if group == "" {
		group = "unknown"
	}
	correlationId := os.Getenv("MON_CORELATIONID")
	if correlationId == "" {
		correlationId = "0"
	}
	taskId := os.Getenv("MESOS_TASK_ID")
	if taskId == "" {
		taskId = "0"
	}

	logger := &DefaultLogger{
		name:          name,
		application:   application,
		group:         group,
		pid:           strconv.Itoa(os.Getpid()),
		correlationId: correlationId,
		taskId:        taskId,
		severity:      INFO,
		severityStreams: map[uint8]severityWriter{
			NOP:     {os.Stdout, "NOP"},
			ALARM:   {os.Stderr, "ALARM"},
			ERROR:   {os.Stderr, "ERROR"},
			STAT:    {os.Stderr, "STAT"},
			INFO:    {os.Stdout, "INFO"},
			EVENT:   {os.Stderr, "EVENT"},
			DEBUG:   {os.Stdout, "DEBUG"},
			UNKNOWN: {os.Stdout, "UNKNOWN"},
		},
	}

	return logger
}

// Debug
// emit using the debug severity level, the only
// optional severity level (see EnableDebug)
func (l *DefaultLogger) Debug(template string, args ...interface{}) {
	if l.severity >= DEBUG {
		l.emit(DEBUG, "", 0, template, args...)
	}
}

// Event
// emit using the event severity level
func (l *DefaultLogger) Event(template string, args ...interface{}) {
	l.emit(EVENT, "", 0, template, args...)
}

// Info
// emit using the info severity level
func (l *DefaultLogger) Info(template string, args ...interface{}) {
	l.emit(INFO, "", 0, template, args...)
}

// Stat
// emit using the stat severity level
func (l *DefaultLogger) Stat(template string, args ...interface{}) {
	l.emit(STAT, "", 0, template, args...)
}

// Error
// emit using the err severity level
func (l *DefaultLogger) Error(template string, args ...interface{}) {
	l.emit(ERROR, "", 0, template, args...)
}

// Alarm
// emit using the alarm severity level
func (l *DefaultLogger) Alarm(template string, args ...interface{}) {
	l.emit(ALARM, "", 0, template, args...)
}

func (l *DefaultLogger) emit(severity uint8, file string, line int, template string, args ...interface{}) {

	// Get caller statistics.
	if file == "" {
		ok := false
		_, file, line, ok = runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
	}

	// Determine the short filename and avoid the func call of strings.SplitAfter.
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}

	file = short
	linx := strconv.Itoa(line)
	fileAndLine := strings.Join([]string{file, linx}, ":")

	// Parse formatted string or use the given template then split into individual lines.
	var lines []string
	if len(args) > 0 {
		formatted := fmt.Sprintf(template, args...)
		lines = strings.Split(formatted, "\n")
	} else {
		lines = strings.Split(template, "\n")
	}

	stream := l.severityStreams[l.severity].writer
	message := strings.Join([]string{
		Marker,
		Version,
		l.severityStreams[l.severity].name,
		l.taskId,
		l.pid,
		l.group,
		l.application,
		l.name,
		l.correlationId,
		fileAndLine}, Separator)

	for _, line := range lines {
		if line == "" {
			continue
		}

		timestamp := time.Now().UTC().Format("2006/01/02 15:04:05.999999999")

		fmt.Fprintln(stream, message+Separator+timestamp+Separator+line)
	}
}
