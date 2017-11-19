package feed

import (
	"crypto/md5"
	"encoding/json"
	"log"
	"strings"

	"github.com/boltdb/bolt"
	strip "github.com/grokify/html-strip-tags-go"
	"github.com/mmcdole/gofeed"
	"github.com/vasilishin/rfeed/slack"
	"github.com/vasilishin/rfeed/store"

	conf "github.com/vasilishin/rfeed/config"
)

var bucket = store.Bucket

// Author feed data
type Author struct {
	Title, Link, Image string
}

// Item can validate items and send to channels
type Item struct {
	Title, Description, Link, Image string
	Author                          *Author
}

func getImage(img *gofeed.Image) (title string, imgURL string) {
	if img != nil {
		title = img.Title
		imgURL = img.URL
	}
	return
}

// NewItem constructor
func NewItem(feed *gofeed.Feed, feedItem *gofeed.Item) Item {
	author := &Author{Link: feed.Link}
	author.Title, author.Image = getImage(feed.Image)

	item := Item{
		Author:      author,
		Title:       strip.StripTags(feedItem.Title),
		Description: strip.StripTags(feedItem.Description),
		Link:        feedItem.Link,
	}
	_, item.Image = getImage(feedItem.Image)
	return item
}

// Save with Item md5 hashable key
func (i *Item) Save() (err error) {
	err = store.Engine.DB.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		jdata, jerr := json.Marshal(i)
		if jerr != nil {
			return err
		}
		return b.Put(i.GetMD5Hash(), jdata)
	})
	return err
}

// Exists Item md5 hashable key in store
func (i *Item) Exists() bool {
	var exists bool
	store.Engine.DB.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		v := b.Get(i.GetMD5Hash())
		exists = v != nil
		return nil
	})
	return exists
}

// Remove by key
func (i *Item) Remove(key []byte) (err error) {
	return store.Engine.DB.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(bucket).Delete(key)
	})
}

// GetMD5Hash return hashable link
func (i *Item) GetMD5Hash() []byte {
	hasher := md5.New()
	hasher.Write([]byte(i.Link))
	return hasher.Sum(nil)
}

// Send message to all channels
func (i *Item) Send() {
	senders := [...]func(){i.ToSlack}
	for _, sender := range senders {
		sender()
	}
}

// ToSlack send message to slack channel
func (i *Item) ToSlack() {
	// create message attachment
	attachment := []*slack.Attachment{
		&slack.Attachment{
			Title:      i.Title,
			Text:       Trim(i.Description),
			TitleLink:  i.Link,
			ImageURL:   i.Image,
			AuthorName: i.Author.Title,
			AuthorLink: i.Author.Link,
			AuthorIcon: i.Author.Image,
		},
	}
	opt := &slack.PostMessageOpt{
		Attachments: attachment,
	}
	if err := slack.Client.SendMessage(i.Title, opt); err != nil {
		log.Fatal(err)
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
		if SkipItem(item) {
			continue
		}
		matchItems = append(matchItems, item)
	}
	return
}

// SkipItem check item categories
func SkipItem(item *gofeed.Item) bool {
	contains := func(categories []string, tag string) bool {
		for _, c := range categories {
			if strings.EqualFold(c, tag) {
				return true
			}
		}
		return false
	}
	for _, tag := range conf.Settings.Tags {
		if contains(item.Categories, tag) {
			return false
		}
	}
	return true
}
