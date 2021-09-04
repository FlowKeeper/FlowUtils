package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Result stores a single result from a item
type Result struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ItemID       primitive.ObjectID
	HostID       primitive.ObjectID
	Type         ReturnType
	CapturedAt   time.Time
	ValueString  string
	ValueNumeric float64
	Error        string
}

type ResultSet struct {
	Results []Result
	Type    ReturnType
}
