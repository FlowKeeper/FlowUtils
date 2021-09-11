package dbtemplate

import (
	"context"
	"time"

	"github.com/FlowKeeper/FlowUtils/v2/models"
	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/logger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

//GetTemplates returns one or multiple template structs for the specified IDs
//If a template id isn't found in the database, no error is caused
//Instead the missing template is omitted from the returned item slice
func GetTemplates(Client *mongo.Database, IDs []primitive.ObjectID) ([]models.Template, error) {
	templates := make([]models.Template, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := Client.Collection("templates").Find(ctx, bson.M{"_id": bson.M{"$in": IDs}})

	if err != nil {
		logger.Error(loggingArea, "Couldn't read template:", err)
		return templates, err
	}

	if err := result.All(ctx, &templates); err != nil {
		logger.Error(loggingArea, "Couldn't decode template array:", err)
	}

	for i := range templates {
		populateTemplateFields(Client, &templates[i])
	}

	return templates, nil
}

//GetAllTemplates returns all templates from the datbase
func GetAllTemplates(Client *mongo.Database) ([]models.Template, error) {
	templates := make([]models.Template, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	result, err := Client.Collection("templates").Find(ctx, bson.M{})

	if err != nil {
		logger.Error(loggingArea, "Couldn't read template:", err)
		return templates, err
	}

	if err := result.All(ctx, &templates); err != nil {
		logger.Error(loggingArea, "Couldn't decode template array:", err)
	}

	for i := range templates {
		populateTemplateFields(Client, &templates[i])
	}

	return templates, nil
}

func populateTemplateFields(Client *mongo.Database, Template *models.Template) {
	//Ensure all arrays != nil
	if Template.ItemIDs == nil {
		Template.ItemIDs = make([]primitive.ObjectID, 0)
	}
	if Template.Items == nil {
		Template.Items = make([]models.Item, 0)
	}
	if Template.TriggerIDs == nil {
		Template.TriggerIDs = make([]primitive.ObjectID, 0)
	}
	if Template.Triggers == nil {
		Template.Triggers = make([]models.Trigger, 0)
	}

	var err error
	if len(Template.ItemIDs) > 0 {
		Template.Items, err = GetItems(Client, Template.ItemIDs)
		if err != nil {
			logger.Error(loggingArea, "Couldn't get items for template", Template.ID, ":", err)
		}
	}
	if len(Template.TriggerIDs) > 0 {
		Template.Triggers, err = GetTriggers(Client, Template.TriggerIDs)
		if err != nil {
			logger.Error(loggingArea, "Couldn't get triggers for template", Template.ID, ":", err)
		}
	}
}
