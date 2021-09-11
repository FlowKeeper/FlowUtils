package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Template struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Name, Description string
	ItemIDs           []primitive.ObjectID
	Items             []Item `bson:"-"`
	TriggerIDs        []primitive.ObjectID
	Triggers          []Trigger `bson:"-"`
}
