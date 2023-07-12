package models

import (
	"fmt"
	"time"

	"github.com/boltdb/bolt"
)

type BoltConfig struct {
	Path string
}

func OpenBolt(config BoltConfig) (*bolt.DB, error) {
	db, err := bolt.Open(config.Path, 0600, &bolt.Options{Timeout: 1 * time.Second})
	if err != nil {
		return nil, fmt.Errorf("open bolt: %w", err)
	}

	return db, nil
}
