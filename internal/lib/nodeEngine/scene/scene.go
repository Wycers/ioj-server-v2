package scene

import (
	"encoding/json"
	"fmt"
)

type V struct {
	Label string      `json:"label"`
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value,omitempty"`
}

type Attribute = map[string]V

type BlockInstance struct {
	ID         int                  `json:"id"`
	Name       string               `json:"name"`
	Title      string               `json:"title"`
	Attributes map[string]Attribute `json:"values"`
}

type Link struct {
	ID         int `json:"id"`
	OriginID   int `json:"originID"`
	OriginSlot int `json:"originSlot"`
	TargetID   int `json:"targetID"`
	TargetSlot int `json:"targetSlot"`
}

type Scene struct {
	Blocks []BlockInstance `json:"blocks"`
	Links  []Link          `json:"links"`
}

func NewScene(jsonStr string) *Scene {
	jsonBytes := []byte(jsonStr)
	scene := new(Scene)
	err := json.Unmarshal(jsonBytes, &scene)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return scene
}
