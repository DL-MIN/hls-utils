package stats

import (
	"bufio"
	"golang.org/x/sys/unix"
	. "hls-utils/logger"
	"hls-utils/terminator"
	"os"
	"strings"
)

type FIFOFile struct {
	Source   string
	DataFunc func(string)
	fd       *os.File
	dataCh   chan string
	errCh    chan error
}

func NewFIFOFile(source string) (f *FIFOFile, err error) {
	Debugf("validate file path %s", source)
	sourceStat, err := os.Stat(source)
	if err != nil && !os.IsNotExist(err) {
		return
	} else if err == nil {
		if sourceStat.IsDir() {
			return nil, ErrLogIsDir(source)
		} else if err = os.Remove(source); err != nil {
			return
		}
	}

	f = &FIFOFile{
		Source:   source,
		DataFunc: func(s string) { Debug(s) },
		dataCh:   make(chan string, 128),
		errCh:    make(chan error),
	}

	Debugf("create fifo file %s", f.Source)
	if err = unix.Mkfifo(f.Source, 0640); err != nil {
		return
	}

	Debugf("open fifo file %s", f.Source)
	if f.fd, err = os.OpenFile(f.Source, os.O_RDWR|unix.O_NONBLOCK, os.ModeNamedPipe); err != nil {
		if err := os.Remove(f.Source); err != nil {
			Warn(err)
		}
		return
	}

	err = nil
	return
}

func (f *FIFOFile) Close() {
	if err := f.fd.Close(); err != nil {
		Warn(err)
	}

	if err := os.Remove(f.Source); err != nil {
		Warn(err)
	}
}

func (f *FIFOFile) ReadPipe() {
	defer f.Close()

	go func() {
		scanner := bufio.NewScanner(f.fd)
		defer func() {
			f.errCh <- scanner.Err()
		}()
		for scanner.Scan() {
			f.dataCh <- strings.TrimSpace(scanner.Text())
		}
	}()

	Debugf("stream from fifo file %s", f.Source)
	for {
		select {
		case <-terminator.Signal:
			return
		case data := <-f.dataCh:
			f.DataFunc(data)
		case err := <-f.errCh:
			Fatal(err)
			return
		}
	}
}
