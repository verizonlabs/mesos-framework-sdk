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

type Logger interface {
	Emit(severity uint8, template string, args ...interface{})
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
	severityStreams map[uint8]severityWriter
}

// Sets required information for our default logger and returns a new instance.
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

	// Intentionally misspelled for historical and backwards compatibility reasons with a few internal tools.
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

// Gets and parses information about the caller.
func (l *DefaultLogger) callerInfo() (string, int) {

	// Get caller file and line number.
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}

	// We only want the filename, not the full path.
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]

	return file, line
}

// Prints out the message to the appropriate stream.
func (l *DefaultLogger) Emit(severity uint8, template string, args ...interface{}) {
	file, line := l.callerInfo()
	fileAndLine := strings.Join([]string{file, strconv.Itoa(line)}, ":")

	// Parse any format specifiers that might be passed in.
	lines := strings.Split(fmt.Sprintf(template, args...), "\n")
	stream := l.severityStreams[severity].writer
	message := strings.Join([]string{
		Marker,
		Version,
		l.severityStreams[severity].name,
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
