package scene

import (
	"encoding/json"
)

type Field struct {
	Name  string `json:"name"`
	Label string `json:"label,omitempty"`
	Type  string `json:"type"`
	Attr  string `json:"attr"`
}

type BlockDefinition struct {
	Name        string  `json:"name"`
	Title       string  `json:"title"`
	Family      string  `json:"family"`
	Description string  `json:"description"`
	Fields      []Field `json:"fields"`
}

func NewBlocksDefinition(jsonStr string) []*BlockDefinition {
	jsonBytes := []byte(jsonStr)
	blocksDefinition := new([]*BlockDefinition)
	err := json.Unmarshal(jsonBytes, &blocksDefinition)
	if err != nil {
		return nil
	}
	return *blocksDefinition
}

func NewBlockDefinition(jsonStr string) *BlockDefinition {
	jsonBytes := []byte(jsonStr)
	blocksDefinition := new(BlockDefinition)
	err := json.Unmarshal(jsonBytes, &blocksDefinition)
	if err != nil {
		return nil
	}
	return blocksDefinition
}
