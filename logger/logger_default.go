package logger

var std = NewLogger()

func Default() *Logger {
    return std
}

func Level() int {
    return std.level
}

func SetLevel(level int) {
    std.level = level
}

func Debug(v ...any) {
    std.logWithLevel(LevelDebug, nil, v...)
}

func Debugf(format string, v ...any) {
    std.logWithLevel(LevelDebug, &format, v...)
}

func Info(v ...any) {
    std.logWithLevel(LevelInfo, nil, v...)
}

func Infof(format string, v ...any) {
    std.logWithLevel(LevelInfo, &format, v...)
}

func Warn(v ...any) {
    std.logWithLevel(LevelWarn, nil, v...)
}

func Warnf(format string, v ...any) {
    std.logWithLevel(LevelWarn, &format, v...)
}

func Fatal(v ...any) {
    std.logWithLevel(LevelFatal, nil, v...)
}

func Fatalf(format string, v ...any) {
    std.logWithLevel(LevelFatal, &format, v...)
}
