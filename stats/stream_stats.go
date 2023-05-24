package stats

import (
    . "hls-utils/logger"
    "hls-utils/terminator"
    "sync"
    "time"
)

const WriteJSONInterval = 1 * time.Second

type StreamStats struct {
    sync.Mutex
    Path string
    data map[string]*StreamStatsData
}

func NewStreamStats(dir string) (s *StreamStats) {
    return &StreamStats{Path: dir, data: make(map[string]*StreamStatsData)}
}

func (s *StreamStats) Add(name string, sequence uint64) {
    s.Lock()
    defer s.Unlock()

    if _, ok := (*s).data[name]; !ok {
        s.newStreamStatsData(name)
    }

    (*s).data[name].hits[sequence]++
}

func (s *StreamStats) close(name string) {
    if _, ok := (*s).data[name]; ok {
        (*s).data[name].Close()
        delete((*s).data, name)
    }
}

func (s *StreamStats) CloseAll() {
    s.Lock()
    defer s.Unlock()

    for i, _ := range (*s).data {
        (*s).data[i].Close()
    }
    (*s).data = nil
}

func (s *StreamStats) Write() {
    for {
        select {
        case <-terminator.Signal:
            s.CloseAll()
            return
        case <-time.After(WriteJSONInterval):
            s.Lock()
            for i, _ := range (*s).data {
                if err := (*s).data[i].write(); err != nil {
                    Info(err)
                    s.close(i)
                }
            }
            s.Unlock()
        }
    }
}
