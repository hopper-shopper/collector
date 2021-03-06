package mongo

import (
	"context"
	"time"

	"github.com/steschwa/hopper-analytics-collector/models"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	BASE_SHARES_COLLECTION = "base-shares"
)

type (
	BaseSharesCollection struct {
		Client *MongoDbClient
	}
)

func (col *BaseSharesCollection) GetCollection() *mongo.Collection {
	return col.Client.Database.Collection(BASE_SHARES_COLLECTION)
}

func (col *BaseSharesCollection) Insert(baseShares models.BaseSharesDocument) error {
	collection := col.GetCollection()

	baseShares.Updated = time.Now()

	_, err := collection.InsertOne(context.Background(), baseShares)

	return err
}
