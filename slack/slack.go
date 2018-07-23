package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/kandziu/rfeed/feed"
)

const postMessageURI = "https://slack.com/api/chat.postMessage"

// Colors for slack Attachment
const (
	Green = "#7CD197"
	Red   = "#F35A00"
)

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

// Send message to channel
func (s *Slack) Send(i *feed.Item) {
	text := fmt.Sprintf("%s: %s", i.Title, i.Link)
	if err := s.SendMessage(text, &PostMessageOpt{}); err != nil {
		log.Fatal(err)
	}
}

// Check client config
func (s *Slack) Check() bool {
	if len(s.token) > 0 && len(s.channel) > 0 {
		return true
	}
	return false
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
