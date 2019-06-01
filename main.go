package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/dmvass/rfeed/telegram"

	conf "github.com/dmvass/rfeed/config"
	"github.com/dmvass/rfeed/store"

	"github.com/dmvass/rfeed/feed"
	"github.com/dmvass/rfeed/slack"
)

// Clients consists from available messangers
var Clients []feed.Messanger

func init() {
	var err error

	// Read settings from config file
	conf.Settings, err = conf.NewSettings("config", ".")
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	// Create connector to Database
	store.Engine, err = store.NewBolt(conf.Settings.Store.Bolt.FilePath)
	if err != nil {
		panic(fmt.Errorf("Fatal error in Database: %s", err))
	}

	// Create message clients
	clients := []feed.Messanger{
		// Create Slack client
		slack.NewClient(conf.Settings.Slack.Token, conf.Settings.Slack.Channel),
		// Create Telegram client
		telegram.NewClient(conf.Settings.Telegram.Token, conf.Settings.Telegram.ChatID),
	}
	for _, c := range clients {
		if c.Check() {
			Clients = append(Clients, c)
		}
	}
	if len(Clients) == 0 {
		panic("You did't specify any message clients.")
	}
}

func main() {
	defer store.Engine.Close()
	// Read feeds every 5 min
	duration := 5 * time.Minute
	wg := new(sync.WaitGroup)
	for _, url := range conf.Settings.Feeds {
		wg.Add(1)
		go observe(url, duration, wg)
	}
	wg.Wait()
}

// Observer for resource
func observe(url string, duration time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(duration)
	for range ticker.C {
		log.Printf("Read from %s resource", url)
		rfeed, err := feed.Read(url)
		if err != nil {
			log.Fatal(err)
		}
		for _, i := range feed.FindItems(rfeed) {
			item := feed.NewItem(i)
			if store.Engine.Exists(item.GetMD5Hash()) {
				continue
			}
			item.Send(&Clients)
			err = store.Engine.Save(item)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
