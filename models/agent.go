package models

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Agent defines the layout of a generic agent stored in the database
type Agent struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Name, Description string
	AgentUUID         uuid.UUID
	Enabled           bool
	LastSeen          time.Time
	OS                AgentOS
	State             AgentState
	TemplateIDs       []primitive.ObjectID
	Templates         []Template `bson:"-"`
	TriggerMappings   []TriggerAssignment
	Endpoint          *url.URL
	ScrapeInterval    int //In seconds
	Scraper           struct {
		UUID uuid.UUID
		Lock time.Time
	}
}

//AgentOS defines on which OS the agent ist running
type AgentOS int

const (
	//Windows AgentOS is set on windows machines
	Windows AgentOS = iota
	//Linux AgentOS is set on linux machines
	Linux
	//Unsupported AgentOS is set if an agent with an unknown OS is registered / loaded
	Unsupported
)

//AgentState defines if the agent is regarded online by the assigned leader
type AgentState int

const (
	//Online is set if the OnlineDetection succeeds
	Online AgentState = iota
	//Offline is set if the OnlineDetection fails
	Offline
)

//AgentosFromString returns the AgentOS iota representation of the specified string
func AgentosFromString(OS string) (AgentOS, error) {
	switch strings.ToLower(OS) {
	case "linux":
		{
			return Linux, nil
		}
	case "windows":
		{
			return Windows, nil
		}
	default:
		{
			return Unsupported, errors.New("unsupported os")
		}
	}
}

//ProblematicTriggers returns all trigger assignments, which are currently in a problematic state
func (a Agent) ProblematicTriggers() []TriggerAssignment {
	problematicTriggers := make([]TriggerAssignment, 0)
	for _, k := range a.TriggerMappings {
		if k.Problematic {
			problematicTriggers = append(problematicTriggers, k)
		}
	}

	return problematicTriggers
}

//GetTrigger returns the trigger struct for the specified ID
func (a Agent) GetTrigger(ID primitive.ObjectID) (Trigger, error) {
	for _, template := range a.Templates {
		for _, trigger := range template.Triggers {
			if trigger.ID == ID {
				return trigger, nil
			}
		}
	}

	return Trigger{}, errors.New("specified trigger wasn't found assigned to agent")
}

//GetItem returns the item struct for the specified ID
func (a Agent) GetItem(ID primitive.ObjectID) (Item, error) {
	for _, template := range a.Templates {
		for _, item := range template.Items {
			if item.ID == ID {
				return item, nil
			}
		}
	}

	return Item{}, errors.New("specified item wasn't found assigned to agent")
}

//GetTriggerMappingByTriggerID returns the TriggerAssignment struct for the specified ID
func (a Agent) GetTriggerMappingByTriggerID(TriggerID primitive.ObjectID) (TriggerAssignment, error) {
	for _, mapping := range a.TriggerMappings {
		if mapping.TriggerID == TriggerID {
			return mapping, nil
		}
	}

	return TriggerAssignment{}, errors.New("specified triggerassignment wasn't found assigned to agent")
}

//GetAllItems returns all items assigned to this agent via templates
//Note that this function already cleans up possibly duplicated items
func (a Agent) GetAllItems() []Item {
	items := make([]Item, 0)
	for _, template := range a.Templates {
		for _, item := range template.Items {
			if !sliceContainsItem(items, item) {
				items = append(items, item)
			}
		}
	}

	return items
}

//GetAllTriggers returns all triggers assigned to this agent via templates
//Note that this function already cleans up possibly duplicated triggers
func (a Agent) GetAllTriggers() []Trigger {
	triggers := make([]Trigger, 0)
	for _, template := range a.Templates {
		for _, trigger := range template.Triggers {
			if !sliceContainsTrigger(triggers, trigger) {
				triggers = append(triggers, trigger)
			}
		}
	}

	return triggers
}

func sliceContainsTrigger(Slice []Trigger, Trigger Trigger) bool {
	for _, k := range Slice {
		if k.ID == Trigger.ID {
			return true
		}
	}

	return false
}

func sliceContainsItem(Slice []Item, Item Item) bool {
	for _, k := range Slice {
		if k.ID == Item.ID {
			return true
		}
	}

	return false
}
