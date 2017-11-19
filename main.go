package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	conf "github.com/vasilishin/rfeed/config"
	"github.com/vasilishin/rfeed/store"

	"github.com/vasilishin/rfeed/feed"
	"github.com/vasilishin/rfeed/slack"
)

func init() {
	var err error
	// Read settings from config file
	conf.Settings, err = conf.NewSettings("config", ".")
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	// Create Slack client
	slack.Client = slack.NewClient(
		conf.Settings.Slack.Token,
		conf.Settings.Slack.Channel,
	)
	// Create connect to Database
	store.Engine, err = store.NewBolt(conf.Settings.Store.Bolt.FilePath)
	if err != nil {
		panic(fmt.Errorf("Fatal error in Database: %s", err))
	}
}

func main() {
	defer store.Engine.DB.Close()
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
			item := feed.NewItem(rfeed, i)
			if item.Exists() {
				continue
			}
			item.Send()
			err = item.Save()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}
