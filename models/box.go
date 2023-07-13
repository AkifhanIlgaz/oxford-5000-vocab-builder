package models

import "github.com/boltdb/bolt"

type BoxService struct {
	DB *bolt.DB
}
