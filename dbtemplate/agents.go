package dbtemplate

import (
	"context"
	"errors"
	"time"

	"github.com/FlowKeeper/FlowUtils/v2/models"
	"github.com/google/uuid"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const loggingArea = "DB"

//GetAgent returns the appropriate agent for the given ID
func GetAgent(Client *mongo.Database, ID primitive.ObjectID) (models.Agent, error) {
	return getAgentByField(Client, "_id", ID)
}

//GetAgent returns the appropriate agent for the given UUID
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

//GetAllAgents returns all agents from the database
func GetAllAgents(Client *mongo.Database) ([]models.Agent, error) {
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
	if Agent.TemplateIDs == nil {
		Agent.TemplateIDs = make([]primitive.ObjectID, 0)
	}
	if Agent.Templates == nil {
		Agent.Templates = make([]models.Template, 0)
	}
	if Agent.TriggerMappings == nil {
		Agent.TriggerMappings = make([]models.TriggerAssignment, 0)
	}

	if len(Agent.TemplateIDs) > 0 {
		var err error
		Agent.Templates, err = GetTemplates(Client, Agent.TemplateIDs)
		if err != nil {
			return err
		}
	}

	//Check if all trigger assignments are still referencing triggers assigned to the agent
	for _, k := range Agent.TriggerMappings {
		if _, err := Agent.GetTrigger(k.TriggerID); err != nil {
			logger.Debug(loggingArea, "Found outdated trigger assignment on agent", Agent.Name, "-> Removing it")
			//ToDo: Remove trigger assignment
		}
	}

	//Check if all assigned triggers have a valid trigger assignment
	for _, trigger := range Agent.GetAllTriggers() {
		if _, err := Agent.GetTriggerMappingByTriggerID(trigger.ID); err != nil {
			logger.Debug(loggingArea, "Found trigger without trigger assignement on agent", Agent.Name, "-> Adding it")
			AddTriggerAssignment(Client, Agent.ID, trigger.ID)
		}
	}

	return nil
}

//AddTriggerAssignment persist a mapping between an agent and a trigger
func AddTriggerAssignment(Client *mongo.Database, AgentID primitive.ObjectID, TriggerID primitive.ObjectID) error {
	return AddTriggerAssignments(Client, AgentID, []primitive.ObjectID{TriggerID})
}

//AddTriggerAssignments persist a mapping between an agent and one or multiple triggers
func AddTriggerAssignments(Client *mongo.Database, AgentID primitive.ObjectID, TriggerIDs []primitive.ObjectID) error {
	newMappings := make([]models.TriggerAssignment, 0)
	for _, k := range TriggerIDs {
		newMappings = append(newMappings, models.TriggerAssignment{
			TriggerID: k,
			Enabled:   true,
			History:   make([]models.TriggerHistoryEntry, 0),
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	result := Client.Collection("agents").FindOneAndUpdate(ctx, bson.M{"_id": AgentID}, bson.M{"$push": bson.M{"triggermappings": bson.M{"$each": newMappings}}})

	if result.Err() != nil {
		logger.Error(loggingArea, "Couldn't add trigger mappings to agent:", result.Err())
	}

	return result.Err()
}
