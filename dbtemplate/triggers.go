package dbtemplate

import (
	"context"
	"time"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/models"
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