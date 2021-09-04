package models

import (
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
