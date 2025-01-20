package types

import (
	"context"
	"io"
	"os"
	"path"
	"regexp"
	"sync"
	"time"

	"github.com/spacecafe/gobox/logger"
)

const (
	PlaylistCheckInterval = time.Second * 30
)

var (
	PlaylistRegex = regexp.MustCompile(`([a-zA-Z0-9._/-]+\\.m3u8)`)
)

type StreamManager struct {
	sync.RWMutex

	ctx context.Context
	wg  sync.WaitGroup

	config  *Config
	streams map[string]*Stream

	// done is a callback function that will be called when the StreamManager is stopped.
	done func()
}

func NewStreamManager(config *Config) *StreamManager {
	return &StreamManager{
		config:  config,
		streams: make(map[string]*Stream),
	}
}

func (r *StreamManager) Start(ctx context.Context, done func()) (err error) {
	r.ctx = ctx
	r.done = done

	go func() {
		<-ctx.Done()
		r.Stop()
	}()

	go r.lifecycleStreams()
	go r.rotateStreamStatistics()

	return
}

func (r *StreamManager) Stop() {
	r.wg.Wait()
	r.done()
}

func (r *StreamManager) AddStream(name string, filename string) (stream *Stream, err error) {
	r.Lock()
	defer r.Unlock()

	logger.Info("adding stream ", name)
	if stream, err = NewStream(name, filename, r.config); err != nil {
		logger.Warn(err)
		return nil, err
	}

	r.streams[name] = stream

	return
}

func (r *StreamManager) GetStream(name string) *Stream {
	r.RLock()
	defer r.RUnlock()

	return r.streams[name]
}

func (r *StreamManager) ListStreams() (streams map[string]*Stream) {
	r.RLock()
	defer r.RUnlock()

	streams = make(map[string]*Stream)
	for name, stream := range r.streams {
		streams[name] = stream
	}

	return
}

func (r *StreamManager) RemoveStream(name string) {
	r.Lock()
	defer r.Unlock()

	logger.Info("removing stream ", name)
	delete(r.streams, name)
}

func (r *StreamManager) GetPlaylist(wr io.Writer, name string) (err error) {
	var stream *Stream
	if stream = r.GetStream(name); stream != nil {
		return stream.GetPlaylist(wr)
	}

	filename := path.Join(r.config.LiveDirectory, name, "index.m3u8")
	if _, err = os.Stat(filename); err == nil {
		if stream, err = r.AddStream(name, filename); err == nil && stream != nil {
			return stream.GetPlaylist(wr)
		}
	}

	return
}

func (r *StreamManager) lifecycleStreams() {
	r.wg.Add(1)
	defer r.wg.Done()

	for {
		select {
		case <-r.ctx.Done():
			return
		case <-time.After(PlaylistCheckInterval):
			for name, stream := range r.ListStreams() {
				if _, err := os.Stat(stream.PlaylistFile); err != nil {
					r.RemoveStream(name)
				}
			}
		}
	}
}

func (r *StreamManager) rotateStreamStatistics() {
	r.wg.Add(1)
	defer r.wg.Done()

	var err error
	for {
		select {
		case <-r.ctx.Done():
			return
		case <-time.After(r.config.StatisticRotationInterval):
			for _, stream := range r.ListStreams() {
				if err = stream.Statistics.Rotate(); err != nil {
					logger.Warn(err)
				}
			}
		}
	}
}
