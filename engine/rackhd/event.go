package rackhd

import (
	"encoding/json"
	"time"
)

type rackhdEvent struct {
	Version   string          `json:"version"`
	Type      string          `json:"type"`
	Action    string          `json:"action"`
	Servrity  string          `json:"servrity"`
	TypeID    string          `json:"typeId"`
	CreatedAt time.Time       `json:"createAt"`
	NodeID    string          `json:"nodeId"`
	Data      json.RawMessage `json:"data"`
}
