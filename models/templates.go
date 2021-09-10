package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Template struct {
	ID         primitive.ObjectID `bson:"_id,omitempty"`
	ItemIDs    []primitive.ObjectID
	Items      []Item `bson:"-"`
	TriggerIDs []primitive.ObjectID
	Triggers   []Trigger `bson:"-"`
}
