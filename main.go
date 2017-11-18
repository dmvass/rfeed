package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/vasilishin/rfeed/store"

	"github.com/spf13/viper"
	"github.com/vasilishin/rfeed/feed"
	"github.com/vasilishin/rfeed/slack"
)

func init() {
	var err error
	// Read settings from config file
	ReadSettings()
	// Create Slack client
	slack.Client = slack.NewClient(
		viper.GetString("slack.token"),
		viper.GetString("slack.channel"),
	)
	// Create connect to Database
	store.Engine, err = store.NewBolt(viper.GetString("db.file"))
	if err != nil {
		panic(fmt.Errorf("Fatal error in Database: %s", err))
	}
}

func main() {
	defer store.Engine.DB.Close()
	// Read feeds every 5 min
	duration := 5 * time.Minute
	wg := new(sync.WaitGroup)
	for _, url := range viper.GetStringSlice("feeds") {
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

// ReadSettings from config file
func ReadSettings() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}
