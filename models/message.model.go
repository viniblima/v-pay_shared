package models

import "encoding/json"

type Message struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}
