package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Trigger struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Name, Description string
	Enabled           bool
	Severity          TriggerSeverity
	DependsOn         []primitive.ObjectID
	Expression        string
}

//TriggerSeverity defines how important a item trigger is
type TriggerSeverity int

const (
	//INFO is the lowest priority trigger
	INFO TriggerSeverity = iota
	//LOW should be used for unimportant triggers
	LOW
	//MEDIUM should be used for important triggers
	MEDIUM
	//HIGH should be used for items which have a high impact in availability,etc.
	HIGH
)
