package types

import (
	"fmt"
	"os"
	"path"
	"sync"
	"time"
)

const (
	CSVHeader = "utc_time,clients\n"
)

// StreamStatistics holds information about connected clients, using two maps to facilitate rotation.
type StreamStatistics struct {
	sync.RWMutex

	// clientsBlue stores client identifiers in one of the two maps used for rotation.
	clientsBlue map[string]struct{}

	// clientsGreen stores client identifiers in the second map used for rotation.
	clientsGreen map[string]struct{}

	// timeline keeps track of the number of clients at different points in time after each rotation.
	timeline map[time.Time]int

	// fd contains the file descriptor.
	fd *os.File

	// useBlue indicates which map (blue or green) is currently active for adding and counting clients.
	useBlue bool
}

// NewStreamStatistics initializes a new StreamStatistics instance,
// setting up both maps and start with clientsBlue.
func NewStreamStatistics(recordDir, name string) (stats *StreamStatistics, err error) {
	stats = &StreamStatistics{
		clientsBlue:  make(map[string]struct{}),
		clientsGreen: make(map[string]struct{}),
		timeline:     make(map[time.Time]int),
		useBlue:      true,
	}

	err = stats.createCSV(recordDir, name)

	return
}

// Add inserts a client identifier into the currently active map.
func (r *StreamStatistics) Add(client string) {
	r.Lock()
	defer r.Unlock()

	r.clientsBlue[client] = struct{}{}
	r.clientsGreen[client] = struct{}{}
}

// Len returns the number of unique clients in the currently active map.
func (r *StreamStatistics) Len() int {
	r.RLock()
	defer r.RUnlock()

	if r.useBlue {
		return len(r.clientsBlue)
	} else {
		return len(r.clientsGreen)
	}
}

// Rotate switches the active map, clearing the previously used one to prepare for new entries.
func (r *StreamStatistics) Rotate() (err error) {
	var currentClients int
	currentTime := time.Now()

	r.Lock()
	r.useBlue = !r.useBlue
	if r.useBlue {
		currentClients = len(r.clientsGreen)
		r.clientsGreen = make(map[string]struct{})
	} else {
		currentClients = len(r.clientsBlue)
		r.clientsBlue = make(map[string]struct{})
	}
	r.timeline[currentTime] = currentClients
	r.Unlock()

	err = r.writeCSVRow(currentTime, currentClients)

	return
}

// Timeline returns a copy of the timeline map containing the history of client counts at different times.
func (r *StreamStatistics) Timeline() map[time.Time]int {
	r.RLock()
	defer r.RUnlock()

	return r.timeline
}

// createCSV creates a CSV file with the given name in the specified directory.
// It appends a timestamp to the filename to ensure uniqueness.
// If the file already exists, it truncates it to zero length before writing.
func (r *StreamStatistics) createCSV(recordDir, name string) (err error) {
	if name != "" {
		filename := path.Join(
			recordDir,
			fmt.Sprintf("%s-%s.csv", name, time.Now().UTC().Format("2006-01-02-15-04-05")),
		)

		err = os.Truncate(path.Clean(filename), 0)
		if err != nil && !os.IsNotExist(err) {
			return err
		}

		r.fd, err = os.OpenFile(path.Clean(filename), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0600)
		if err != nil {
			return
		}
		_, err = r.fd.WriteString(CSVHeader)
	}

	return
}

// writeCSVRow writes a row to the CSV file with the given time and number of clients.
// It formats the time in UTC as "YYYY-MM-DD HH:MM:SS" and appends it to the file.
func (r *StreamStatistics) writeCSVRow(time time.Time, clients int) (err error) {
	if r.fd != nil {
		_, err = fmt.Fprintf(r.fd, "%s,%d\n", time.UTC().Format("2006-01-02 15:04:05"), clients)
	}

	return
}
