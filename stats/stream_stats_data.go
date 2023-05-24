package stats

import (
    "errors"
    "fmt"
    "github.com/gogo/protobuf/sortkeys"
    "golang.org/x/sys/unix"
    . "hls-utils/logger"
    "os"
    "path"
)

type StreamStatsData struct {
    Close        func()
    fileJSON     *os.File
    filePlaylist string
    hits         map[uint64]uint64
    lastSequence uint64
    lastRead     int
}

func (s *StreamStats) newStreamStatsData(name string) {
    fileJSON := path.Join(s.Path, name+".json")
    file, err := os.OpenFile(fileJSON, os.O_CREATE|os.O_WRONLY, 0640)
    if err != nil {
        Warn(err)
        return
    }

    (*s).data[name] = &StreamStatsData{
        Close: func() {
            if err := file.Close(); err != nil {
                Warn(err)
            }
            if err := os.Remove(fileJSON); err != nil {
                Warn(err)
            }
        },
        fileJSON:     file,
        filePlaylist: path.Join(s.Path, name+".m3u8"),
        hits:         make(map[uint64]uint64),
    }
}

func (s *StreamStatsData) getSortedKeys() (keys []uint64) {
    keys = make([]uint64, 0, len((*s).hits))
    for item := range (*s).hits {
        keys = append(keys, item)
    }
    sortkeys.Uint64s(keys)
    return
}

func (s *StreamStatsData) getMinOfN(n int) (min uint64) {
    keys := (*s).getSortedKeys()
    if s.lastSequence == keys[len(keys)-1] {
        if s.lastRead >= n {
            min = 0
            return
        }
        s.lastRead++
    } else {
        s.lastRead = 0
    }

    s.lastSequence = keys[len(keys)-1]
    min = (*s).hits[s.lastSequence]

    if len(keys) > n {
        for i := 0; i < n; i++ {
            min = Min(min, (*s).hits[keys[len(keys)-1-i]])
        }
    }
    return
}

func (s *StreamStatsData) write() error {
    subscribers := s.getMinOfN(3)
    _, err := os.Stat(s.filePlaylist)
    if err != nil {
        return ErrPlaylistNotExist(s.filePlaylist)
    }

    go func() {
        if err := unix.Flock(int(s.fileJSON.Fd()), unix.LOCK_EX|unix.LOCK_NB); err == nil {
            defer unix.Flock(int(s.fileJSON.Fd()), unix.LOCK_UN)

            if err := s.fileJSON.Truncate(0); err != nil {
                Warn(err)
                return
            }

            if _, err := s.fileJSON.Seek(0, 0); err != nil {
                Warn(err)
                return
            }

            if _, err := fmt.Fprintf(s.fileJSON, `{"subscribers":%d}`, subscribers); err != nil {
                Warn(err)
            }
        } else if !errors.Is(err, unix.EAGAIN) {
            Warn(err)
        }
    }()

    return nil
}
