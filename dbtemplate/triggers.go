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

func GetTrigger(Client *mongo.Database, ID primitive.ObjectID) (models.Trigger, error) {
	triggers, err := GetTriggers(Client, []primitive.ObjectID{ID})
	if err != nil {
		return models.Trigger{}, err
	}

	return triggers[0], nil
}

func GetTriggerByName(Client *mongo.Database, Name string) (models.Trigger, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result := Client.Collection("triggers").FindOne(ctx, bson.M{"name": Name})

	if result.Err() != nil {
		if errors.Is(result.Err(), mongo.ErrNoDocuments) {
			return models.Trigger{}, result.Err()
		}

		logger.Error(loggingArea, "Couldn't read item:", result.Err())
		return models.Trigger{}, result.Err()
	}

	var trigger models.Trigger
	if err := result.Decode(&trigger); err != nil {
		logger.Error(loggingArea, "Couldn't decode item:", err)
	}

	return trigger, nil
}

func GetTriggers(Client *mongo.Database, IDs []primitive.ObjectID) ([]models.Trigger, error) {
	triggers := make([]models.Trigger, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := Client.Collection("triggers").Find(ctx, bson.M{"_id": bson.M{"$in": IDs}})

	if err != nil {
		logger.Error(loggingArea, "Couldn't read items:", err)
		return triggers, err
	}

	if err := result.All(ctx, &triggers); err != nil {
		logger.Error(loggingArea, "Couldn't decode trigger array:", err)
	}

	return triggers, nil
}
