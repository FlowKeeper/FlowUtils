package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Item struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Name, Description string
	Returns           ReturnType
	Unit              string
	Interval          int //Interval is set in seconds
	Command           string
	CheckOn           AgentOS //Execute this item only on this os

}

//ReturnType defines which type of information is returned by the check
type ReturnType int

const (
	//Numeric is set if the check returns a number
	Numeric ReturnType = iota
	//Text is set if the check returns text
	Text
)

//Result stores a single result from a item
type Result struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ItemID       primitive.ObjectID
	HostID       primitive.ObjectID
	CapturedAt   time.Time
	ValueString  string
	ValueNumeric float64
	Error        string
}
