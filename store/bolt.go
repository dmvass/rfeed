package store

import (
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var (
	// Bucket name for feeds
	Bucket = []byte("feeds")
	// Engine global client
	Engine *Bolt
)

// Bolt implements store.Engine with boltdb
type Bolt struct {
	DB *bolt.DB
}

// NewBolt makes persitent boltdb based store
func NewBolt(dbFile string) (*Bolt, error) {
	log.Printf("bolt (persitent) store, %s", dbFile)
	store := Bolt{}
	db, err := bolt.Open(dbFile, 0600, &bolt.Options{Timeout: 1 * time.Second})
	db.Update(func(tx *bolt.Tx) error {
		_, e := tx.CreateBucketIfNotExists(Bucket)
		return e
	})
	store.DB = db
	return &store, err
}
