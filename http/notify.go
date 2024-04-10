package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/spf13/viper"
	. "hls-utils/logger"
)

func notify(name string) {
	data := map[string]string{"name": name}

	jsonData, err := json.Marshal(data)
	if err != nil {
		Warn(err)
		return
	}

	notifyURL, err := url.Parse(viper.GetString("notify"))
	if err != nil {
		Warn(err)
		return
	}

	req, err := http.NewRequest("POST", notifyURL.String(), bytes.NewBuffer(jsonData))
	if err != nil {
		Warn(err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		Warn(err)
		return
	}
	defer resp.Body.Close()
}
