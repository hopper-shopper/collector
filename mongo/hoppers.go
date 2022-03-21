package mongo

import (
	"context"

	"github.com/steschwa/hopper-analytics-collector/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	HOPPERS_COLLECTION = "hoppers"
)

type (
	HoppersCollection struct {
		Connection *mongo.Client
	}
)

func (col *HoppersCollection) GetCollection() *mongo.Collection {
	return GetCollection(col.Connection, HOPPERS_COLLECTION)
}

func (col *HoppersCollection) Upsert(hopper models.HopperDocument) error {
	collection := col.GetCollection()

	upsert := true
	_, err := collection.ReplaceOne(
		context.Background(),
		bson.D{{Key: "tokenId", Value: hopper.TokenId}},
		hopper,
		&options.ReplaceOptions{Upsert: &upsert},
	)

	return err
}
