package stats

import (
    "github.com/spf13/viper"
    . "hls-utils/logger"
    "hls-utils/terminator"
    "path"
    "regexp"
    "strconv"
)

func Run() {
    stats := NewStreamStats(path.Clean(viper.GetString("stats.data")))
    regex, err := regexp.Compile(viper.GetString("stats.regex"))
    if err != nil {
        Fatal(err)
    }
    nameIdx := regex.SubexpIndex("name")
    sequenceIdx := regex.SubexpIndex("sequence")

    fifo, err := NewFIFOFile(path.Clean(viper.GetString("stats.log")))
    if err != nil {
        Fatal(err)
    }
    fifo.DataFunc = func(data string) {
        matches := regex.FindStringSubmatch(data)
        if matches != nil {
            if sequence, err := strconv.ParseUint(matches[sequenceIdx], 10, 64); err == nil {
                stats.Add(matches[nameIdx], sequence)
            }
        }
    }

    terminator.WaitGroup.Add(2)
    go func() {
        fifo.ReadPipe()
        terminator.WaitGroup.Done()
    }()

    go func() {
        stats.Write()
        terminator.WaitGroup.Done()
    }()
}
