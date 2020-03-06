package feed

import (
	"crypto/md5"
	"strings"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"

	conf "github.com/dmvass/rfeed/config"
)

// Messanger interface for clients
type Messanger interface {
	Send(i *Item)
	Check() bool
}

// Item can validate items and send to channels
type Item struct {
	Title, Link string
	Origin      *gofeed.Item
}

// NewItem constructor
func NewItem(feedItem *gofeed.Item) *Item {
	item := &Item{
		Title:  strip.StripTags(feedItem.Title),
		Link:   feedItem.Link,
		Origin: feedItem,
	}
	return item
}

// GetMD5Hash return hashable link
func (i *Item) GetMD5Hash() []byte {
	hasher := md5.New()
	hasher.Write([]byte(i.Link))
	return hasher.Sum(nil)
}

// Send message to all channels
func (i *Item) Send(clients *[]Messanger) {
	for _, client := range *clients {
		go client.Send(i)
	}
}

// Read resource
func Read(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return feed, err
}

// FindItems return filtered items
func FindItems(feed *gofeed.Feed) (matchItems []*gofeed.Item) {
	for _, item := range feed.Items {
		if skipItem(item) {
			continue
		}
		matchItems = append(matchItems, item)
	}
	return
}

func skipItem(item *gofeed.Item) bool {
	contains := func(categories []string, tag string) bool {
		for _, c := range categories {
			if strings.EqualFold(c, tag) {
				return true
			}
		}
		return false
	}
	if conf.Settings.Tags == nil {
		return false
	}
	for _, tag := range conf.Settings.Tags {
		if contains(item.Categories, tag) {
			return false
		}
	}
	return true
}
