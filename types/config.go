package types

import (
	"errors"
	"flag"
	"net/url"
	"os"
	"time"

	"github.com/spacecafe/gobox/httpserver"
	"github.com/spacecafe/gobox/logger"
	"github.com/spf13/viper"
)

var (
	DefaultLogLevel                  = "debug"
	DefaultLiveDirectory             = "/tmp"
	DefaultRecordDirectory           = "/tmp"
	DefaultStatisticRotationInterval = time.Second * 10

	ErrInvalidLiveDirectory             = errors.New("data directory must exist")
	ErrInvalidRecordDirectory           = errors.New("record directory must exist")
	ErrInvalidNotificationEndpoint      = errors.New("notification endpoint is not a valid URL")
	ErrInvalidStreams                   = errors.New("streams must not be nil")
	ErrInvalidTrackLabels               = errors.New("track labels must not be nil")
	ErrInvalidStatisticRotationInterval = errors.New("statistic rotation interval must be greater than 0")
	ErrNoPassword                       = errors.New("password of a stream cannot be empty")
)

// Config defines the essential parameters for serving this application.
type Config struct {
	// LogLevel specifies the level of logging to be used.
	LogLevel string `json:"log_level" yaml:"log_level" mapstructure:"log_level"`

	// LiveDirectory is the path to the directory where live HLS data is stored.
	LiveDirectory string `json:"live_dir" yaml:"live_dir" mapstructure:"live_dir"`

	// RecordDirectory is the path to the directory where recorded HLS data is stored.
	RecordDirectory string `json:"record_dir" yaml:"record_dir" mapstructure:"record_dir"`

	// NotificationEndpoint is the URL endpoint where notifications will be sent.
	NotificationEndpoint string `json:"notification_endpoint" yaml:"notification_endpoint" mapstructure:"notification_endpoint"`

	// Streams is a map of stream names to their respective password.
	Streams map[string]string `json:"streams" yaml:"streams" mapstructure:"streams"`

	// TrackLabels is a map of track names to their labels.
	TrackLabels map[string]string `json:"track_labels" yaml:"track_labels" mapstructure:"track_labels"`

	StatisticRotationInterval time.Duration `json:"statistic_rotation_interval" yaml:"statistic_rotation_interval" mapstructure:"statistic_rotation_interval"`

	// HTTPServer holds the configuration for the embedded HTTP server.
	HTTPServer *httpserver.Config `json:"http_server" yaml:"http_server" mapstructure:"http_server"`
}

// NewConfig creates and returns a new Config having default values from given configuration file.
func NewConfig() *Config {
	config := &Config{
		LogLevel:                  DefaultLogLevel,
		LiveDirectory:             DefaultLiveDirectory,
		RecordDirectory:           DefaultRecordDirectory,
		Streams:                   make(map[string]string),
		TrackLabels:               make(map[string]string),
		StatisticRotationInterval: DefaultStatisticRotationInterval,
		HTTPServer:                httpserver.NewConfig(logger.Default()),
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	configPath := flag.String("config", "", "Path to config.yaml")
	flag.Parse()
	if *configPath != "" {
		viper.SetConfigFile(*configPath)
	} else {
		viper.AddConfigPath("/etc/hls-utils/")
		viper.AddConfigPath(".")
	}

	err := viper.ReadInConfig()
	if err != nil {
		logger.Fatal(err)
	}

	err = viper.Unmarshal(config)
	if err != nil {
		logger.Fatal(err)
	}

	if err = config.Validate(); err != nil {
		logger.Fatal(err)
	}

	return config
}

// Validate ensures the all necessary configurations are filled and within valid confines.
// Any misconfiguration results in well-defined standardized errors.
func (r *Config) Validate() error {
	var err error
	var info os.FileInfo
	if err = logger.ParseLevel(r.LogLevel); err != nil {
		return err
	}

	if info, err = os.Stat(r.LiveDirectory); os.IsNotExist(err) || !info.IsDir() {
		return ErrInvalidLiveDirectory
	}

	if info, err = os.Stat(r.RecordDirectory); os.IsNotExist(err) || !info.IsDir() {
		return ErrInvalidRecordDirectory
	}

	if r.NotificationEndpoint != "" {
		if _, err = url.Parse(r.NotificationEndpoint); err != nil {
			return ErrInvalidNotificationEndpoint
		}
	}

	if r.Streams == nil {
		return ErrInvalidStreams
	}
	for i := range r.Streams {
		if r.Streams[i] == "" {
			return ErrNoPassword
		}
	}

	if r.TrackLabels == nil {
		return ErrInvalidTrackLabels
	}

	if r.StatisticRotationInterval <= 0 {
		return ErrInvalidStatisticRotationInterval
	}

	if err = r.HTTPServer.Validate(); err != nil {
		return err
	}

	return nil
}
