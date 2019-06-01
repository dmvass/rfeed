package store

import (
	"os"
	"reflect"
	"testing"

	"github.com/dmvass/rfeed/feed"
)

// test data
var tables = []struct {
	Title, Link string
}{
	{"Test 1", "http://example.com/feed/1/"},
	{"Test 2", "http://example.com/feed/2/"},
	{"Test 3", "http://example.com/feed/3/"},
	{"Test 4", "http://example.com/feed/4/"},
	{"Test 5", "http://example.com/feed/5/"},
}

func TestCreateBolt(t *testing.T) {
	// create test boltdb file
	b, err := NewBolt("./test.bd")
	defer os.Remove("./test.bd")
	if err != nil {
		t.Errorf("Fatal error in db file: %s", err)
	}
	err = b.Close()
	if err != nil {
		t.Errorf("Fatal error when the connection is closed: %s", err)
	}
}

func TestSaveAndLoadBolt(t *testing.T) {
	// create test boltdb file
	b, err := NewBolt("./test.bd")
	defer os.Remove("./test.bd")
	if err != nil {
		t.Errorf("Fatal error in db file: %s", err)
	}

	for _, table := range tables {
		// create test item
		item := &feed.Item{Title: table.Title, Link: table.Link}
		// save item to boltdb
		err = b.Save(item)
		if err != nil {
			t.Errorf("Item %v can't be save to boltdb: %s", item, err)
		}
		// load item from boltdb
		loadItem, err := b.Load(item.GetMD5Hash())
		if err != nil {
			t.Errorf("Item %v can't be load from boltdb: %s", item, err)
		}
		// load item must be identical to created item
		if !reflect.DeepEqual(item, loadItem) {
			t.Errorf("Load item was incorrect, got: %v, want: %v", item, loadItem)
		}
		// load item must not be identical to created item
		item.Link = "http://example.com/feed/changed/"
		if reflect.DeepEqual(item, loadItem) {
			t.Errorf("Load item was incorrect, got: %v, want: %v", item, loadItem)
		}

	}

	// load not exists item from boltdb
	loadItem, err := b.Load([]byte("not exists key"))
	if err == nil {
		t.Errorf("Load incorect item %v from boltdb", loadItem)
	}

}

func TestExistsAndRemoveKeyBolt(t *testing.T) {
	// create test boltdb file
	b, err := NewBolt("./test.bd")
	defer os.Remove("./test.bd")
	if err != nil {
		t.Errorf("Fatal error in db file: %s", err)
	}
	for _, table := range tables {
		// create test item
		item := &feed.Item{Title: table.Title, Link: table.Link}
		// save item to boltdb
		err = b.Save(item)
		if err != nil {
			t.Errorf("Item %v can't be save to boltdb: %s", item, err)
		}
		// item hash key must be in boltdb
		if !b.Exists(item.GetMD5Hash()) {
			t.Errorf("Item %v not exists in boltdb, got: false, want: true", item)
		}
		// remove item from boltdb
		err = b.Remove(item.GetMD5Hash())
		if err != nil {
			t.Errorf("Item %v can't be deleted from boltdb: %s", item, err)
		}
		// item hash key must not be in boltdb
		if b.Exists(item.GetMD5Hash()) {
			t.Errorf("Item %v exists in boltdb, got: true, want: false", item)
		}
	}
}
