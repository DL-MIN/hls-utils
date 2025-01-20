package types

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Notification struct {
	Name string `json:"name"`
}

func (r *Notification) Send(endpoint string) (err error) {
	payload, err := json.Marshal(r)
	if err != nil {
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	return
}
