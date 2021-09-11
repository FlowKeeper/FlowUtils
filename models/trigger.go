package models

import (
	"time"

	"gitlab.cloud.spuda.net/Wieneo/golangutils/v2/stringHelper"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Trigger specifies the layout of a generic trigger stored in the database
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

//TriggerAssignment is used to map a trigger (specified via the TriggerID) to an agent
//TriggerAssignments are automatically created by the dbtemplate package upon retrieving an agent
type TriggerAssignment struct {
	Enabled     bool
	TriggerID   primitive.ObjectID
	Problematic bool
	Error       string
	History     []TriggerHistoryEntry
}

//TriggerHistoryEntry is used to store when a trigger became problematic / unproblematic for a given TriggerAssignment / TriggerMapping
type TriggerHistoryEntry struct {
	Time        time.Time
	Problematic bool
}

//HasError returns true if the Error string is set to something other than ""
func (t TriggerAssignment) HasError() bool {
	return !stringHelper.IsEmpty(t.Error)
}
