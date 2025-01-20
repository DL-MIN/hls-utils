package types

import (
	"io"
	"os"
	"strings"
	"text/template"

	"github.com/google/uuid"
)

type Stream struct {
	config           *Config
	Statistics       *StreamStatistics
	PlaylistFile     string
	playlistTemplate *template.Template
}

func NewStream(name, filename string, config *Config) (stream *Stream, err error) {
	stream = &Stream{
		config:       config,
		PlaylistFile: filename,
	}
	stream.Statistics, err = NewStreamStatistics(config.RecordDirectory, name)
	if err != nil {
		return nil, err
	}
	err = stream.CreatePlaylistTemplate()

	return
}

func (r *Stream) GetPlaylist(wr io.Writer) (err error) {
	return r.playlistTemplate.Execute(wr, struct {
		ClientID string
	}{
		ClientID: uuid.NewString(),
	})
}

func (r *Stream) CreatePlaylistTemplate() (err error) {
	playlistRaw, err := os.ReadFile(r.PlaylistFile)
	if err != nil {
		return
	}
	playlist := string(playlistRaw)

	for name, label := range r.config.TrackLabels {
		playlist = strings.ReplaceAll(playlist, name, label)
	}

	playlist = PlaylistRegex.ReplaceAllString(playlist, "{{.ClientID}}/$1")

	r.playlistTemplate, err = template.New(r.PlaylistFile).Parse(playlist)

	return
}
