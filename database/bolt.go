package database

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

type BoltConfig struct {
	FileName string
}

func OpenBolt(config BoltConfig) (*bolt.DB, error) {
	db, err := bolt.Open(config.FileName, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("open bolt: %w", err)
	}

	return db, nil
}
