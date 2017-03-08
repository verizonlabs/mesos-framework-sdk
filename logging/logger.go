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

func (l *DefaultLogger) Emit(severity uint8, template string, args ...interface{}) {

	// Get caller statistics.
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}

	// Determine the short filename.
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
