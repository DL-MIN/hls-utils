package stats

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path"

	"github.com/gogo/protobuf/sortkeys"
	"golang.org/x/sys/unix"
	. "hls-utils/logger"
)

// AnalyzeNSequences defines the amount of last the sequences to be analyzed
const AnalyzeNSequences = 5

// ResponseJSON is used as a template to store statistics in the web server's path
type ResponseJSON struct {
	Status int    `json:"status"`
	Code   string `json:"code"`
	Title  string `json:"title"`
	Data   struct {
		Subscribers struct {
			Current uint64             `json:"current"`
			History *map[uint64]uint64 `json:"history"`
		} `json:"subscribers"`
	} `json:"data"`
}

// newStreamStatsDataJSON creates and fills a response structure
func (s *StreamStatsData) newStreamStatsDataJSON() (sj *ResponseJSON) {
	sj = &ResponseJSON{
		Status: http.StatusOK,
		Code:   "load_hls_statistics",
		Title:  "Successfully load HLS statistics",
	}
	sj.Data.Subscribers.Current = s.getMaxOfN(AnalyzeNSequences)
	sj.Data.Subscribers.History = &s.hits

	return
}

// StreamStatsData is used in StreamStats to collect the amount of viewer of a streaming endpoint
type StreamStatsData struct {
	// close releases the file descriptor and removes the JSON file
	close func()

	// fileJSON contains the file descriptor of the JSON file
	fileJSON *os.File

	// filePlaylist is the path to the video playlist file
	filePlaylist string

	// hits counts viewers per video sequence
	hits map[uint64]uint64

	// lastSequence is set to last to detect zero viewers
	// It is set to the last key in hits on every write event.
	lastSequence uint64

	// lastRead is set to detect zero viewers.
	// It is incremented by 1 if lastSequence is equal to the last key in hits.
	lastRead int
}

// newStreamStatsData creates and open a new JSON file. A new StreamStatsData is returned.
func (s *StreamStats) newStreamStatsData(name string) {
	file, err := os.OpenFile(path.Clean(path.Join(s.Path, name+".json")), os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		Warn(err)
		return
	}

	(*s).data[name] = &StreamStatsData{
		close: func() {
			if err := file.Close(); err != nil {
				Warn(err)
			}
			if err := os.Remove(file.Name()); err != nil {
				Warn(err)
			}
		},
		fileJSON:     file,
		filePlaylist: path.Clean(path.Join(s.Path, name+".m3u8")),
		hits:         make(map[uint64]uint64),
	}
}

// getSortedKeys returns an ascending ordered array of video sequences
func (s *StreamStatsData) getSortedKeys() (keys []uint64) {
	keys = make([]uint64, 0, len((*s).hits))
	for item := range (*s).hits {
		keys = append(keys, item)
	}
	sortkeys.Uint64s(keys)
	return
}

// getMaxOfN determines the current amount of viewers by returning the smallest amount of the last n sequences
func (s *StreamStatsData) getMaxOfN(n int) (max uint64) {
	keys := (*s).getSortedKeys()
	if s.lastSequence == keys[len(keys)-1] {
		if s.lastRead >= n {
			max = 0
			return
		}
		s.lastRead++
	} else {
		s.lastRead = 0
	}

	s.lastSequence = keys[len(keys)-1]
	max = (*s).hits[s.lastSequence]

	if len(keys) > n {
		for i := 0; i < n; i++ {
			max = Max(max, (*s).hits[keys[len(keys)-1-i]])
		}
	}
	return
}

// write stores the current amount of viewers to the JSON file
func (s *StreamStatsData) write() error {
	responseJSON, err := json.Marshal(s.newStreamStatsDataJSON())
	if err != nil {
		return err
	}

	_, err = os.Stat(s.filePlaylist)
	if err != nil {
		return ErrPlaylistNotExist(s.filePlaylist)
	}

	go func() {
		if err := unix.Flock(int(s.fileJSON.Fd()), unix.LOCK_EX|unix.LOCK_NB); err == nil {
			defer func() {
				if err := unix.Flock(int(s.fileJSON.Fd()), unix.LOCK_UN); err != nil {
					Warn(err)
				}
			}()

			if err := s.fileJSON.Truncate(0); err != nil {
				Warn(err)
				return
			}

			if _, err := s.fileJSON.Seek(0, 0); err != nil {
				Warn(err)
				return
			}

			if _, err := fmt.Fprintf(s.fileJSON, `%s`, responseJSON); err != nil {
				Warn(err)
			}
		} else if !errors.Is(err, unix.EAGAIN) {
			Warn(err)
		}
	}()

	return nil
}
