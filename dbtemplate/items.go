package dbtemplate

import (
	"context"
	"errors"
	"time"

	"github.com/FlowKeeper/FlowUtils/v2/models"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetItems(Client *mongo.Database, IDs []primitive.ObjectID) ([]models.Item, error) {
	items := make([]models.Item, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := Client.Collection("items").Find(ctx, bson.M{"_id": bson.M{"$in": IDs}})

	if err != nil {
		logger.Error(loggingArea, "Couldn't read items:", err)
		return items, err
	}

	if err := result.All(ctx, &items); err != nil {
		logger.Error(loggingArea, "Couldn't decode item array:", err)
	}

	return items, nil
}

func GetItemByName(Client *mongo.Database, Name string) (models.Item, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result := Client.Collection("items").FindOne(ctx, bson.M{"name": Name})

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return models.Item{}, result.Err()
		}

		logger.Error(loggingArea, "Couldn't read item:", result.Err())
		return models.Item{}, result.Err()
	}

	var item models.Item
	if err := result.Decode(&item); err != nil {
		logger.Error(loggingArea, "Couldn't decode item:", err)
	}

	return item, nil
}
