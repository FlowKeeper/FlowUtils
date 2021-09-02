package models

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Agent struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	AgentID   uuid.UUID
	ScraperID uuid.UUID
	Enabled   bool
	LastSeen  time.Time
	OS        AgentOS
	State     AgentState
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
