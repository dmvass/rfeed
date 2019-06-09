package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmvass/rfeed/pool"
	"github.com/dmvass/rfeed/telegram"

	conf "github.com/dmvass/rfeed/config"
	"github.com/dmvass/rfeed/store"

	"github.com/dmvass/rfeed/feed"
	"github.com/dmvass/rfeed/slack"
)

var configPath = flag.String("config", "./config.yml", "config file path")

// Clients consists from available messangers
var Clients []feed.Messanger

func init() {
	var err error

	flag.Parse()

	// Read settings from config file
	conf.Settings, err = conf.NewSettings(*configPath)
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
	duration := time.Duration(conf.Settings.Interval) * time.Second

	sigs := make(chan os.Signal, 1)
	defer close(sigs)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	workerPool := pool.New(len(conf.Settings.Feeds))
	workerPool.Run()

	go func(p *pool.Pool) {
		<-sigs
		p.Close()
	}(workerPool)

	go observe(workerPool, duration)

	workerPool.Wait()
}

// Observer for resource
func observe(p *pool.Pool, duration time.Duration) {
	ticker := time.NewTicker(duration)
	for range ticker.C {
		for _, url := range conf.Settings.Feeds {

			job := func(url string) func() {
				return func() {
					log.Printf("Read from %s resource", url)
					rfeed, err := feed.Read(url)
					if err != nil {
						log.Print(err)
						return
					}
					for _, i := range feed.FindItems(rfeed) {
						item := feed.NewItem(i)
						if store.Engine.Exists(item.GetMD5Hash()) {
							continue
						}
						item.Send(&Clients)
						err = store.Engine.Save(item)
						if err != nil {
							log.Print(err)
							continue
						}
					}
				}
			}(url)

			p.Submit(job)
		}
	}
}
