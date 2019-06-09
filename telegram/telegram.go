package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/dmvass/rfeed/feed"
)

// Telegram API
const (
	API               = "https://api.telegram.org/bot%s/%s"
	SendMessageMethod = "sendMessage"
)

// Message type
type Message struct {
	ChatID int64  `json:"chat_id"`
	Text   string `json:"text"`
}

// Telegram client
type Telegram struct {
	token  string
	chatID int64
}

// NewClient create new telegram client
func NewClient(token string, chatID int64) *Telegram {
	telegram := &Telegram{
		token:  token,
		chatID: chatID,
	}
	return telegram
}

// Check client config
func (t *Telegram) Check() bool {
	if len(t.token) > 0 && t.chatID > 0 {
		return true
	}
	return false
}

// Send message
func (t *Telegram) Send(i *feed.Item) {
	text := fmt.Sprintf("%s: %s", i.Title, i.Link)
	if err := t.SendMessage(text); err != nil {
		log.Fatal(err)
	}
}

// SendMessage to telegram chat by id
func (t *Telegram) SendMessage(text string) error {
	log.Printf("Send to telegram chat %d: %s", t.chatID, text)
	message := &Message{
		ChatID: t.chatID,
		Text:   text,
	}
	data, err := json.Marshal(message)
	if err != nil {
		return err
	}
	// build request to telegram API
	req, err := t.buildRequest(data, SendMessageMethod)
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

// build telegram request
func (t *Telegram) buildRequest(data []byte, method string) (*http.Request, error) {
	endpoint := fmt.Sprintf(API, t.token, method)
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}
