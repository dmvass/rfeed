package slack

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const postMessageURI = "https://slack.com/api/chat.postMessage"

// Colors for slack Attachment
const (
	Green = "#7CD197"
	Red   = "#F35A00"
)

// Client global slack client
var Client *Slack

// Slack client
type Slack struct {
	token, channel string
}

// NewClient slack client builder
func NewClient(token, channel string) *Slack {
	slack := &Slack{
		token:   "Bearer " + token,
		channel: channel,
	}
	return slack
}

// SendMessage send message to slack channel
func (s *Slack) SendMessage(text string, opt *PostMessageOpt) error {
	log.Printf("Send to slack: %s", text)
	postMessage := &PostMessage{
		Channel:        s.channel,
		Text:           text,
		PostMessageOpt: *opt,
	}
	data, err := json.Marshal(postMessage)
	if err != nil {
		return err
	}
	// build request to slack API with Authorization header
	req, err := s.buildRequest(data)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Fatal(string(body))
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (s *Slack) buildRequest(data []byte) (*http.Request, error) {
	req, err := http.NewRequest("POST", postMessageURI, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", s.token)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}