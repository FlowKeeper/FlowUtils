package models

import (
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Agent struct {
	ID                primitive.ObjectID `bson:"_id,omitempty"`
	Name, Description string
	AgentUUID         uuid.UUID
	Enabled           bool
	LastSeen          time.Time
	OS                AgentOS
	State             AgentState
	Items             []primitive.ObjectID
	ItemsResolved     []Item `bson:"-"`
	Triggers          []TriggerAssignment
	Endpoint          *url.URL
	Scraper           struct {
		UUID uuid.UUID
		Lock time.Time
	}
}

//AgentOS defines on which OS the agent ist running
type AgentOS int

const (
	//Windows is set on windows machines
	Windows AgentOS = iota
	//Linux is set on linux machines
	Linux
	Unsupported
)

//HostState defines if the host is regarded online by the assigned leader
type AgentState int

const (
	//Online is set if the OnlineDetection succeeds
	Online AgentState = iota
	//Offline is set if the OnlineDetection fails
	Offline
)

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
	for _, k := range a.Triggers {
		if k.Problematic {
			problematicTriggers = append(problematicTriggers, k)
		}
	}

	return problematicTriggers
}
