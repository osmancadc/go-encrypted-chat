package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

type Logger struct {
	level   string
	output  io.Writer
	format  string
	context map[string]interface{}
}

func NewLogger(level string) *Logger {
	return &Logger{
		level:   level,
		output:  os.Stdout,
		format:  "text",
		context: make(map[string]interface{}),
	}
}

func (l *Logger) SetOutput(output io.Writer) {
	l.output = output
}

func (l *Logger) SetFormat(format string) {
	l.format = format
}

func (l *Logger) WithContext(ctx context.Context) *Logger {
	newLogger := *l
	newLogger.context = extractContextFields(ctx)
	return &newLogger
}

func (l *Logger) Log(message string, fields ...map[string]interface{}) {
	l.logln("CHAT", message, fields...)
}

func (l *Logger) Logf(format string, v ...interface{}) {
	l.log("CHAT", fmt.Sprintf(format, v...))
}

func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	l.logln("DEBUG", message, fields...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.log("DEBUG", fmt.Sprintf(format, v...))
}

func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	l.logln("INFO", message, fields...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.log("INFO", fmt.Sprintf(format, v...))
}

func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	l.logln("WARN", message, fields...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.log("WARN", fmt.Sprintf(format, v...))
}

func (l *Logger) Error(message string, fields ...map[string]interface{}) {
	l.logln("ERROR", message, fields...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.log("ERROR", fmt.Sprintf(format, v...))
}

func (l *Logger) Fatal(message string, fields ...map[string]interface{}) {
	l.logln("FATAL", message, fields...)
	os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.log("FATAL", fmt.Sprintf(format, v...))
	os.Exit(1)
}

func (l *Logger) log(level string, message string, fields ...map[string]interface{}) {
	if l.shouldLog(level) {
		timestamp := time.Now().Format(time.RFC3339)
		logMessage := l.formatMessage(timestamp, level, message, fields...)
		fmt.Fprint(l.output, logMessage)
	}
}

func (l *Logger) logln(level string, message string, fields ...map[string]interface{}) {
	if l.shouldLog(level) {
		timestamp := time.Now().Format(time.RFC3339)
		logMessage := l.formatMessage(timestamp, level, message, fields...)
		fmt.Fprintln(l.output, logMessage)
	}
}

func (l *Logger) formatMessage(timestamp, level, message string, fields ...map[string]interface{}) string {

	logEntry := fmt.Sprintf("%s [%s] %s", timestamp, level, message)

	if l.format == "json" {
		logEntry = fmt.Sprintf(`{"timestamp":"%s","level":"%s","message":"%s","fields":%v}`, timestamp, level, message, fields)
	} else if len(fields) > 0 {
		for _, fieldMap := range fields {
			if fieldMap != nil {
				fieldData, err := json.Marshal(fieldMap)
				if err != nil {
					log.Printf("Failed to serialize fields: %v", err)
				} else {
					logEntry = fmt.Sprintf("%s %s", logEntry, string(fieldData))
				}
			}
		}

	}

	return logEntry
}

func (l *Logger) shouldLog(level string) bool {
	levels := map[string]int{"DEBUG": 1, "INFO": 2, "WARN": 3, "ERROR": 4, "FATAL": 5, "CHAT": 6}
	return levels[level] >= levels[l.level]
}

func extractContextFields(_ context.Context) map[string]interface{} {
	fields := make(map[string]interface{})
	return fields
}
