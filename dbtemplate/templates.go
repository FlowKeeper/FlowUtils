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

	for i, k := range templates {
		if len(k.ItemIDs) > 0 {
			templates[i].Items, err = GetItems(Client, k.ItemIDs)
			if err != nil {
				logger.Error(loggingArea, "Couldn't get items for template", k.ID, ":", err)
			}
		}
		if len(k.TriggerIDs) > 0 {
			templates[i].Triggers, err = GetTriggers(Client, k.TriggerIDs)
			if err != nil {
				logger.Error(loggingArea, "Couldn't get triggers for template", k.ID, ":", err)
			}
		}
	}

	return templates, nil
}
