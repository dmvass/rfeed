package store

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	"github.com/vasilishin/rfeed/feed"
)

var (
	// Bucket name for feeds
	Bucket = []byte("feeds")
	// Engine global client
	Engine *Bolt
)

// DB Errors
var (
	ErrLoadRejected = fmt.Errorf("message expired or deleted")
)

// Bolt implements store.Engine with boltdb
type Bolt struct {
	db *bolt.DB
}

// NewBolt makes persitent boltdb based store
func NewBolt(dbFile string) (*Bolt, error) {
	log.Printf("bolt (persitent) store, %s", dbFile)
	store := Bolt{}
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, err
	}
	db.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists(Bucket)
		return e
	})
	store.db = db
	return &store, err
}

// Close db public method
func (b *Bolt) Close() error {
	err := b.db.Close()
	return err
}

// Save with Item hashable key in store
func (b *Bolt) Save(i *feed.Item) (err error) {
	err = b.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket(Bucket)
		jdata, jerr := json.Marshal(i)
		if jerr != nil {
			return err
		}
		return b.Put(i.GetMD5Hash(), jdata)
	})
	return err
}

// Load by key, removes on first access, checks expire
func (b *Bolt) Load(key []byte) (i *feed.Item, err error) {
	err = b.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(Bucket)
		val := bucket.Get(key)
		if val == nil {
			log.Printf("Key not found %s", key)
			return ErrLoadRejected
		}
		i = &feed.Item{}
		return json.Unmarshal(val, i)
	})

	return i, err
}

// Exists Item md5 hashable key in store
func (b *Bolt) Exists(key []byte) bool {
	var exists bool
	b.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(Bucket)
		v := b.Get(key)
		exists = v != nil
		return nil
	})
	return exists
}

// Remove by key from store
func (b *Bolt) Remove(key []byte) (err error) {
	return b.db.Update(func(tx *bolt.Tx) error {
		return tx.Bucket(Bucket).Delete(key)
	})
}
