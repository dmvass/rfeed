package slack

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dmvass/rfeed/feed"
)

const postMessageURI = "https://slack.com/api/chat.postMessage"

// PostMessage post message
type PostMessage struct {
	Channel string `json:"channel"`
	Text    string `json:"text"`
	PostMessageOpt
}

// PostMessageOpt option type for `chat.postMessage` api
type PostMessageOpt struct {
	AsUser      bool          `json:"as_user"`
	Username    string        `json:"username"`
	Parse       string        `json:"parse"`
	LinkNames   string        `json:"link_names"`
	Attachments []*Attachment `json:"attachments"`
	UnfurlLinks string        `json:"unfurl_links"`
	UnfurlMedia string        `json:"unfurl_media"`
	IconURL     string        `json:"icon_url"`
	IconEmoji   string        `json:"icon_emoji"`
}

// attachField it is possible to create more richly-formatted
// messages using Attachments. https://api.slack.com/docs/attachments
type attachField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

// Attachment for message
type Attachment struct {
	Color         string         `json:"color,omitempty"`
	Fallback      string         `json:"fallback"`
	AuthorName    string         `json:"author_name,omitempty"`
	AuthorSubname string         `json:"author_subname,omitempty"`
	AuthorLink    string         `json:"author_link,omitempty"`
	AuthorIcon    string         `json:"author_icon,omitempty"`
	Title         string         `json:"title,omitempty"`
	TitleLink     string         `json:"title_link,omitempty"`
	Pretext       string         `json:"pretext,omitempty"`
	Text          string         `json:"text"`
	ImageURL      string         `json:"image_url,omitempty"`
	ThumbURL      string         `json:"thumb_url,omitempty"`
	Footer        string         `json:"footer,omitempty"`
	FooterIcon    string         `json:"footer_icon,omitempty"`
	TimeStamp     int64          `json:"ts,omitempty"`
	Fields        []*attachField `json:"fields,omitempty"`
	MarkdownIn    []string       `json:"mrkdwn_in,omitempty"`
}

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
