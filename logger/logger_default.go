package logger

var std = NewLogger()

func Default() *Logger {
	return std
}

var (
	Level    = (*std).Level
	SetLevel = (*std).SetLevel
	Debug    = (*std).Debug
	Debugf   = (*std).Debugf
	Info     = (*std).Info
	Infof    = (*std).Infof
	Warn     = (*std).Warn
	Warnf    = (*std).Warnf
	Fatal    = (*std).Fatal
	Fatalf   = (*std).Fatalf
)
