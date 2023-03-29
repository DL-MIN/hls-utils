/*******************************************************************************
 * parse access.log to receive the current number of HLS subscribers
 *
 * logformat: [UNIX TIMESTAMP (or msec)] [FILE]
 * e.g.:      1639399770.294 /srv/domain.tld/hls/stream-name-12.ts
 *
 * @author Lars Thoms
 * @date   2023-03-29
 ******************************************************************************/

package main

import (
    "bufio"
    "flag"
    "log"
    "os"
    "regexp"
    "strconv"
    "time"
)

/*******************************************************************************
 * distribution struct
 ******************************************************************************/

type distribution struct {
    dist map[uint64]uint64
    max  uint64
}

func (d *distribution) Init() {
    d.dist = make(map[uint64]uint64)
}

func (d *distribution) Count(i uint64) {
    d.dist[i]++
    if d.max < i {
        d.max = i
    }
}

func (d *distribution) Get(i uint64) uint64 {
    return d.dist[i]
}

/*******************************************************************************
 * functions
 ******************************************************************************/

/**
 * @brief      Parse access.log and generate a distribution
 *
 * @param      path    Path to access.log
 * @param      stream  Name of the stream
 * @param      dist    Distribution
 */
func ParseLog(path *string, stream *string, dist *distribution) {

    // open access.log
    file, err := os.Open(*path)
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()

    // prepare filter
    time := time.Now().Unix() / 100
    regex := regexp.MustCompile(
        "^" + strconv.FormatInt(time, 10) + "[\\d]{2}\\.[\\d]+.*?" + regexp.QuoteMeta(*stream) + "_(?:src|\\d+p)-([\\d]+).ts$")

    // iterate trough logfile
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        matches := regex.FindStringSubmatch(scanner.Text())
        if matches != nil {
            i, _ := strconv.ParseUint(matches[1], 10, 64)
            dist.Count(i)
        } else {
            continue
        }
    }

    if err := scanner.Err(); err != nil {
        log.Fatal(err)
    }
}

/**
 * @brief      Generate stats.json with the current number of subscribers
 *
 * @param      dist      Distribution
 * @param      path      Path to stats.json
 * @param      segments  Number of stream segments to analyse
 */
func GenerateStat(dist *distribution, path *string, segments uint64) {
    var subscribers uint64

    for i := uint64(0); i < segments; i++ {
        if s := dist.Get(dist.max - i); subscribers < s {
            subscribers = s
        }
    }

    // write stats as JSON to file
    err := os.WriteFile(*path, []byte("{\"subscribers\":"+strconv.FormatUint(subscribers, 10)+"}"), 0644)
    if err != nil {
        log.Fatal(err)
    }
}

/**
 * @brief      HLS Subscriber Stats
 */
func main() {
    // CLI arguments
    logfilePtr := flag.String("input", "access.log", "Path to access.log")
    statsfilePtr := flag.String("output", "stream.json", "Path to stats.json")
    streamPtr := flag.String("name", "stream", "Name of the stream")
    segmentsPtr := flag.Uint64("segments", 3, "Number of stream segments to analyse")
    intervalPtr := flag.Int("interval", 10, "Pause between parsing operations in seconds")
    flag.Parse()

    // parse log all n seconds
    for {
        distPtr := new(distribution)
        distPtr.Init()
        ParseLog(logfilePtr, streamPtr, distPtr)
        GenerateStat(distPtr, statsfilePtr, *segmentsPtr)

        time.Sleep(time.Duration(*intervalPtr) * time.Second)
    }
}
