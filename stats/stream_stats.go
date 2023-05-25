package stats

import (
	. "hls-utils/logger"
	"hls-utils/terminator"
	"sync"
	"time"
)

// WriteJSONInterval defines the amount of time to wait between writing JSON files
const WriteJSONInterval = 1 * time.Second

// StreamStats collects the amount of viewers per streaming endpoint
type StreamStats struct {
	sync.Mutex

	// Path refers to the directory of video playlist files
	Path string

	// data contains statistics per streaming endpoint
	data map[string]*StreamStatsData
}

// NewStreamStats returns a new StreamStats with initialized fields
func NewStreamStats(dir string) (s *StreamStats) {
	return &StreamStats{Path: dir, data: make(map[string]*StreamStatsData)}
}

// Add a new hit of a video sequence to the specified streaming endpoint
func (s *StreamStats) Add(name string, sequence uint64) {
	s.Lock()
	defer s.Unlock()

	if _, ok := (*s).data[name]; !ok {
		s.newStreamStatsData(name)
	}

	(*s).data[name].hits[sequence]++
}

// close releases the file descriptor to the JSON file and removes the streaming endpoint from data
func (s *StreamStats) close(name string) {
	if _, ok := (*s).data[name]; ok {
		(*s).data[name].close()
		delete((*s).data, name)
	}
}

// CloseAll calls close to all streaming endpoints
func (s *StreamStats) CloseAll() {
	s.Lock()
	defer s.Unlock()

	for i := range (*s).data {
		(*s).data[i].close()
	}
	(*s).data = nil
}

// Write stores every WriteJSONInterval the current amount of viewers per streaming endpoint to a JSON file
func (s *StreamStats) Write() {
	for {
		select {
		case <-terminator.Signal:
			s.CloseAll()
			return
		case <-time.After(WriteJSONInterval):
			s.Lock()
			for i := range (*s).data {
				if err := (*s).data[i].write(); err != nil {
					Info(err)
					s.close(i)
				}
			}
			s.Unlock()
		}
	}
}
