package setup

import (
	"fmt"

	"github.com/AkifhanIlgaz/vocab-builder/database"
	"github.com/boltdb/bolt"
	"go.mongodb.org/mongo-driver/mongo"
)

type databases struct {
	Mongo *mongo.Client
	Bolt  *bolt.DB
}

// TODO: Use MongoDb instead of Bolt for box

func Databases(config *config) (*databases, error) {
	mongo, err := database.OpenMongo(config.Mongo)
	if err != nil {
		return nil, fmt.Errorf("setup databases | mongo: %w", err)
	}

	bolt, err := database.OpenBolt(config.Bolt)
	if err != nil {
		return nil, fmt.Errorf("setup databases | bolt: %w", err)
	}

	return &databases{
		Mongo: mongo,
		Bolt:  bolt,
	}, nil
}
