package models

import "go.mongodb.org/mongo-driver/bson/primitive"

//Template specifies the layout of a generic template stored in the database

type Template struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Name, Description string
	ItemIDs           []primitive.ObjectID
	Items             []Item `bson:"-"`
	TriggerIDs        []primitive.ObjectID
	Triggers          []Trigger `bson:"-"`
}
