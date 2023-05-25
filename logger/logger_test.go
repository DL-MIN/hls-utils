package logger

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"io"
	"log"
	"os"
	"os/exec"
	"testing"
)

func TestLogger_Print(t *testing.T) {
	type fields struct {
		level int
	}
	type args struct {
		format string
		v      []any
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{"debug one string", fields{LevelDebug}, args{"%s", []any{"TEST"}}, "[\x1b[0;37mDEBUG\x1b[0m] TEST\n"},
		{"debug two strings", fields{LevelDebug}, args{"%s%s", []any{"TE", "ST"}}, "[\x1b[0;37mDEBUG\x1b[0m] TEST\n"},
		{"debug one string and integer", fields{LevelDebug}, args{"%s%d", []any{"TEST", 123}}, "[\x1b[0;37mDEBUG\x1b[0m] TEST123\n"},
		{"info one string", fields{LevelInfo}, args{"%s", []any{"TEST"}}, "[\x1b[0;32mINFO\x1b[0m] TEST\n"},
		{"info two strings", fields{LevelInfo}, args{"%s%s", []any{"TE", "ST"}}, "[\x1b[0;32mINFO\x1b[0m] TEST\n"},
		{"info one string and integer", fields{LevelInfo}, args{"%s%d", []any{"TEST", 123}}, "[\x1b[0;32mINFO\x1b[0m] TEST123\n"},
		{"warn one string", fields{LevelWarn}, args{"%s", []any{"TEST"}}, "[\x1b[0;33mWARNING\x1b[0m] TEST\n"},
		{"warn two strings", fields{LevelWarn}, args{"%s%s", []any{"TE", "ST"}}, "[\x1b[0;33mWARNING\x1b[0m] TEST\n"},
		{"warn one string and integer", fields{LevelWarn}, args{"%s%d", []any{"TEST", 123}}, "[\x1b[0;33mWARNING\x1b[0m] TEST123\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			Default().loggerList[tt.fields.level].SetOutput(&buf)
			Default().loggerList[tt.fields.level].SetFlags(0)

			switch tt.fields.level {
			case LevelDebug:
				Debug(tt.args.v...)
			case LevelInfo:
				Info(tt.args.v...)
			case LevelWarn:
				Warn(tt.args.v...)
			}

			assert.Equal(t, tt.want, buf.String())
			buf.Reset()

			switch tt.fields.level {
			case LevelDebug:
				Debugf(tt.args.format, tt.args.v...)
			case LevelInfo:
				Infof(tt.args.format, tt.args.v...)
			case LevelWarn:
				Warnf(tt.args.format, tt.args.v...)
			}

			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestLogger_Fatal(t *testing.T) {
	if os.Getenv("CRASHTEST") == "1" {
		Default().loggerList[LevelFatal].SetFlags(0)
		Fatalf("%s%d", "TEST", 123)
		return
	}

	cmd := exec.Command(os.Args[0], append(os.Args[1:], "-test.run=TestLogger_Fatal")...)
	cmd.Env = append(os.Environ(), "CRASHTEST=1", "GOCOVERDIR=/tmp")
	bufReader, _ := cmd.StderrPipe()
	err := cmd.Start()
	bufOut, _ := io.ReadAll(bufReader)
	err = cmd.Wait()

	e, ok := err.(*exec.ExitError)
	assert.Equal(t, true, ok)
	assert.Equal(t, false, e.Success())

	want := "[\x1b[0;31mFATAL\x1b[0m] TEST123\n"
	assert.EqualValues(t, want, string(bufOut))
}

func TestLogger_Level(t *testing.T) {
	type fields struct {
		level int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{"debug level", fields{LevelDebug}, LevelDebug},
		{"info level", fields{LevelInfo}, LevelInfo},
		{"warn level", fields{LevelWarn}, LevelWarn},
		{"fatal level", fields{LevelFatal}, LevelFatal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &Logger{
				level: tt.fields.level,
			}
			if got := l.Level(); got != tt.want {
				t.Errorf("Level() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLogger_SetLevel(t *testing.T) {
	type fields struct {
		level int
	}
	type args struct {
		level int
	}
	tests := []struct {
		name  string
		args  args
		wants fields
	}{
		{"debug level", args{LevelDebug}, fields{LevelDebug}},
		{"info level", args{LevelInfo}, fields{LevelInfo}},
		{"warn level", args{LevelWarn}, fields{LevelWarn}},
		{"fatal level", args{LevelFatal}, fields{LevelFatal}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetLevel(tt.args.level)
		})
		assert.Equal(t, tt.wants.level, Level())
	}
}

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name string
		want *Logger
	}{
		{name: "valid logger", want: &Logger{level: 0, loggerList: [4]*log.Logger{log.Default(), log.Default(), log.Default()}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewLogger()
			assert.Equal(t, tt.want.level, got.level)
			for i := range got.loggerList {
				assert.IsType(t, tt.want.loggerList[i], got.loggerList[i])
			}
		})
	}
}

func TestLogger_logWithLevel(t *testing.T) {
	type args struct {
		level  int
		format *string
		v      []any
	}
	tests := []struct {
		name     string
		args     args
		setLevel int
		want     string
	}{
		{"debug level", args{LevelDebug, nil, []any{"TEST"}}, LevelDebug, "[\x1b[0;37mDEBUG\x1b[0m] TEST\n"},
		{"info level", args{LevelInfo, nil, []any{"TEST"}}, LevelInfo, "[\x1b[0;32mINFO\x1b[0m] TEST\n"},
		{"warn level", args{LevelWarn, nil, []any{"TEST"}}, LevelWarn, "[\x1b[0;33mWARNING\x1b[0m] TEST\n"},
		{"debug at warn level", args{LevelDebug, nil, []any{"TEST"}}, LevelWarn, ""},
		{"info at warn level", args{LevelInfo, nil, []any{"TEST"}}, LevelWarn, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			SetLevel(tt.setLevel)
			Default().loggerList[tt.setLevel].SetOutput(&buf)
			Default().loggerList[tt.setLevel].SetFlags(0)

			Default().logWithLevel(tt.args.level, tt.args.format, tt.args.v...)

			assert.Equal(t, tt.want, buf.String())
		})
	}
}
