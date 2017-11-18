package slack

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
