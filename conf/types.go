package conf

import (
	"encoding/json"
	"strings"
)

const (
	ReportPendingStatus = iota
	ReportRunningStatus
	ReportSucceededStatus
	ReportFailedStatus
	ReportUnknowStatus
)

type ReportStatus uint

func (r ReportStatus) String() string {
	var s string
	switch r {
	case ReportPendingStatus:
		s = "pending"
	case ReportRunningStatus:
		s = "running"
	case ReportSucceededStatus:
		s = "succeeded"
	case ReportFailedStatus:
		s = "failed"
	default:
		s = "unknow"
	}
	return s
}

type Config struct {
	Scenarios []*Scenario `json:"scenarios"`
}

type Scenario struct {
	Name   string          `json:"name"`
	Nodes  []Node          `json:"nodes"`
	Label  []string        `json:"label"`
	Engine string          `json:"engine"`
	Data   json.RawMessage `json:"data"`
}

type ScenarioReport struct {
	Name   string
	Status ReportStatus
	Data   []byte
}

type Node struct {
	Name        string   `json:"name"`
	Identifiers []string `json:"identifiers"`
}

func (n Node) String() (s string) {
	if len(n.Name) != 0 {
		s = n.Name
	} else {
		s = strings.Join(n.Identifiers, ",")
	}

	return s
}
