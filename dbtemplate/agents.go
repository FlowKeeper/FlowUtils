package dbtemplate

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"gitlab.cloud.spuda.net/flowkeeper/flowutils/v2/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const loggingArea = "DB"

func GetAgent(Client *mongo.Database, ID primitive.ObjectID) (models.Agent, error) {
	return getAgentByField(Client, "_id", ID)
}

func GetAgentByUUID(Client *mongo.Database, UUID uuid.UUID) (models.Agent, error) {
	return getAgentByField(Client, "agentuuid", UUID)
}

func getAgentByField(Client *mongo.Database, Field string, Value interface{}) (models.Agent, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result := Client.Collection("agents").FindOne(ctx, bson.M{Field: Value})

	if result.Err() != nil {
		if !errors.Is(result.Err(), mongo.ErrNoDocuments) {
			logger.Error(loggingArea, "Couldn't fetch agent from db:", result.Err())
		}

		return models.Agent{}, result.Err()
	}

	var agent models.Agent
	if err := result.Decode(&agent); err != nil {
		logger.Error(loggingArea, "Couldn't decode agent:", err)
		return models.Agent{}, err
	}

	//Fix if array is nil
	if agent.Items == nil {
		agent.Items = make([]primitive.ObjectID, 0)
	}
	if agent.ItemsResolved == nil {
		agent.ItemsResolved = make([]models.Item, 0)
	}
	if agent.Triggers == nil {
		agent.Triggers = make([]models.TriggerAssignment, 0)
	}

	if len(agent.Items) > 0 {
		var err error
		agent.ItemsResolved, err = GetItems(Client, agent.Items)
		if err != nil {
			return agent, err
		}
	}

	for i, k := range agent.Triggers {
		var err error
		agent.Triggers[i].Trigger, err = GetTrigger(Client, k.TriggerID)
		if err != nil {
			logger.Error("Couldn't resolve trigger", k.TriggerID, ":", err)
			return agent, err
		}
	}

	return agent, nil
}
