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

//GetTrigger gets the appropriate trigger which matches the specified ID
func GetTrigger(Client *mongo.Database, ID primitive.ObjectID) (models.Trigger, error) {
	triggers, err := GetTriggers(Client, []primitive.ObjectID{ID})
	if err != nil {
		return models.Trigger{}, err
	}

	return triggers[0], nil
}

//GetTriggerByName gets the appropriate trigger which matches the specified Name
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

//GetTriggers returns one or multiple trigger structs for the specified IDs
//If a trigger id isn't found in the database, no error is caused
//Instead the missing trigger is omitted from the returned item slice
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
