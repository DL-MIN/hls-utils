package logger

// std is the default logger
var std = NewLogger()

// Default returns the default logger
func Default() *Logger {
	return std
}

// Aliases of logger functions
var (
	Level    = std.Level
	SetLevel = std.SetLevel
	Debug    = std.Debug
	Debugf   = std.Debugf
	Info     = std.Info
	Infof    = std.Infof
	Warn     = std.Warn
	Warnf    = std.Warnf
	Fatal    = std.Fatal
	Fatalf   = std.Fatalf
)
