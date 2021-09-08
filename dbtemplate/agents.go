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

	if err := populateAgentFields(Client, &agent); err != nil {
		return models.Agent{}, err
	}

	return agent, nil
}

func GetAgents(Client *mongo.Database) ([]models.Agent, error) {
	var agents []models.Agent

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	result, err := Client.Collection("agents").Find(ctx, bson.M{})

	if err != nil {
		logger.Error(loggingArea, "Couldn't fetch agents from db:", err)
		return agents, err
	}

	err = result.All(ctx, &agents)

	if err != nil {
		logger.Error(loggingArea, "Couldn't decode agents:", err)
	}

	if agents == nil {
		agents = make([]models.Agent, 0)
	}

	for i := range agents {
		if err := populateAgentFields(Client, &agents[i]); err != nil {
			return agents, err
		}
	}

	return agents, nil
}

func populateAgentFields(Client *mongo.Database, Agent *models.Agent) error {
	//Fix if array is nil
	if Agent.Items == nil {
		Agent.Items = make([]primitive.ObjectID, 0)
	}
	if Agent.ItemsResolved == nil {
		Agent.ItemsResolved = make([]models.Item, 0)
	}
	if Agent.Triggers == nil {
		Agent.Triggers = make([]models.TriggerAssignment, 0)
	}

	if len(Agent.Items) > 0 {
		var err error
		Agent.ItemsResolved, err = GetItems(Client, Agent.Items)
		if err != nil {
			return err
		}
	}

	for i, k := range Agent.Triggers {
		var err error
		Agent.Triggers[i].Trigger, err = GetTrigger(Client, k.TriggerID)
		if err != nil {
			logger.Error("Couldn't resolve trigger", k.TriggerID, ":", err)
			return err
		}
	}

	return nil
}
