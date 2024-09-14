package valueobj

import (
	"encoding/json"
)

type TaskDetail struct {
	Content string `json:"content"`
}

func (t TaskDetail) Marshal() string {
	raw, _ := json.Marshal(t)
	return string(raw)
}

func (t TaskDetail) Unmarshal(raw string) error {
	return json.Unmarshal([]byte(raw), &t)
}
