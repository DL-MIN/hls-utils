package logger

import (
	"fmt"
	"log"
	"os"
)

const (
	// LevelDebug is indented only for debugging errors and misbehaviors.
	// Do not use this in a production environment to improve performance.
	LevelDebug = 0 + iota

	// LevelInfo is useful for application usage analysis
	LevelInfo

	// LevelWarn is used for non-critical situations that require immediate abort (panic).
	LevelWarn

	// LevelFatal indicates a critical situation requiring immediate abort (panic).
	// The application is terminated with a SIGTERM and returns 1 as exit code.
	LevelFatal
)

// A Logger allows messages with different levels/priorities to be sent to stderr/stdout
type Logger struct {
	level      int
	loggerList [4]*log.Logger
}

// NewLogger returns a new logger at the debug level
func NewLogger() *Logger {
	return &Logger{
		level: LevelDebug,
		loggerList: [4]*log.Logger{
			log.New(os.Stdout, "[\u001B[0;37mDEBUG\u001B[0m] ", log.Ldate|log.Ltime|log.Lshortfile),
			log.New(os.Stdout, "[\u001B[0;32mINFO\u001B[0m] ", log.Ldate|log.Ltime|log.Lshortfile),
			log.New(os.Stderr, "[\u001B[0;33mWARNING\u001B[0m] ", log.Ldate|log.Ltime|log.Lshortfile),
			log.New(os.Stderr, "[\u001B[0;31mFATAL\u001B[0m] ", log.Ldate|log.Ltime|log.Lshortfile),
		},
	}
}

// Level returns the current log level
func (l *Logger) Level() int {
	return l.level
}

// SetLevel sets the current log level to the desired value. Avoid misbehavior by using package defined constants.
func (l *Logger) SetLevel(level int) {
	l.level = level
}

// logWithLevel writes message to stderr/stdout in accordance with its level
func (l *Logger) logWithLevel(level int, format *string, v ...any) {
	if l.level > level {
		return
	}

	var err error

	if format == nil {
		err = l.loggerList[level].Output(3, fmt.Sprint(v...))
	} else {
		err = l.loggerList[level].Output(3, fmt.Sprintf(*format, v...))
	}

	if err != nil {
		fmt.Println(err)
	}

	if level == LevelFatal {
		os.Exit(1)
	}
}

// Debug writes debug level messages
func (l *Logger) Debug(v ...any) {
	l.logWithLevel(LevelDebug, nil, v...)
}

// Debugf writes debug level messages using formatted string
func (l *Logger) Debugf(format string, v ...any) {
	l.logWithLevel(LevelDebug, &format, v...)
}

// Info writes info level messages
func (l *Logger) Info(v ...any) {
	l.logWithLevel(LevelInfo, nil, v...)
}

// Infof writes info level messages using formatted string
func (l *Logger) Infof(format string, v ...any) {
	l.logWithLevel(LevelInfo, &format, v...)
}

// Warn writes warn level messages
func (l *Logger) Warn(v ...any) {
	l.logWithLevel(LevelWarn, nil, v...)
}

// Warnf writes warn level messages using formatted string
func (l *Logger) Warnf(format string, v ...any) {
	l.logWithLevel(LevelWarn, &format, v...)
}

// Fatal writes fatal level messages and terminates the application
func (l *Logger) Fatal(v ...any) {
	l.logWithLevel(LevelFatal, nil, v...)
}

// Fatalf writes fatal level messages using formatted string and terminates the application
func (l *Logger) Fatalf(format string, v ...any) {
	l.logWithLevel(LevelFatal, &format, v...)
}
