package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/infinity-oj/server-v2/internal/lib/engine/scene"
)

type Node struct {
	Id         int                    `json:"id"`
	Type       string                 `json:"type"`
	Properties map[string]interface{} `json:"properties"`
	Inputs     []interface{}          `json:"inputs"`
	Outputs    []interface{}          `json:"outputs"`
}

type Output struct {
	Name string      `json:"name"`
	Type string      `json:"type"`
	Link interface{} `json:"link"`
}

type GraphX struct {
	Nodes []Node        `json:"nodes"`
	Links []interface{} `json:"links"`
}

func NewGraphByFile(filename string) (*Graph, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	return NewGraphByDefinition(string(data))
}

func NewGraphByScene(bs []*scene.BlockDefinition, s *scene.Scene) (*Graph, error) {
	graph := New()

	blockMap := make(map[string]*scene.BlockDefinition)

	for _, b := range bs {
		blockMap[b.Name] = b
	}

	for _, v := range s.Blocks {

		inputCounts := 0
		outputCounts := 0
		if b, ok := blockMap[v.Name]; ok {
			for _, field := range b.Fields {
				if field.Attr == "input" {
					inputCounts++
				}
				if field.Attr == "output" {
					outputCounts++
				}
			}
		}

		var inputs []int
		for i := 0; i < inputCounts; i++ {
			for _, link := range s.Links {
				if link.TargetID == v.ID && link.TargetSlot == i {
					inputs = append(inputs, link.ID)
				}
			}
		}

		var outputs [][]int
		for i := 0; i < outputCounts; i++ {
			outputs = append(outputs, []int{})
		}

		block := graph.AddBlock(v.ID, v.Name, nil, inputs, outputs)

		for k, attr := range v.Attributes["property"] {
			if attr.Value == nil {
				continue
			}
			block.setProperty(k, attr.Value)
		}
	}

	for _, v := range s.Links {
		graph.AddLink(
			v.ID,
			v.OriginID,
			v.OriginSlot,
			v.TargetID,
			v.TargetSlot,
		)
	}

	return graph, nil
}

func NewGraphByDefinition(definition string) (*Graph, error) {
	var dataObject GraphX
	err := json.Unmarshal([]byte(definition), &dataObject)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	graph := New()

	for _, v := range dataObject.Nodes {

		var inputs []int
		for _, input := range v.Inputs {
			if mp, ok := input.(map[string]interface{}); ok {
				if link, ok := mp["link"]; ok {
					if link == nil {
						continue
					}
					inputs = append(inputs, int(link.(float64)))
				}
			}
		}
		//for _, link := range links {
		//	if link == nil {
		//		continue
		//	}
		//	inputs = append(inputs, link.(int))
		//}
		var outputs [][]int
		for _, output := range v.Outputs {
			if mp, ok := output.(map[string]interface{}); ok {
				if links, ok := mp["links"]; ok {
					if links == nil {
						outputs = append(outputs, []int{0})
						continue
					}
					var tmp []int
					for _, link := range links.([]interface{}) {
						if link == nil {
							continue
						}
						tmp = append(tmp, int(link.(float64)))
					}
					outputs = append(outputs, tmp)
				}
			}
		}

		block := graph.AddBlock(v.Id, v.Type, nil, inputs, outputs)

		for k, v := range v.Properties {
			if v, ok := v.(string); ok {
				block.setProperty(k, v)
			} else {
				fmt.Println("==> value not string:", k, v)
			}
		}
	}

	for _, v := range dataObject.Links {
		if arr, ok := v.([]interface{}); ok {
			graph.AddLink(
				int(arr[0].(float64)),
				int(arr[1].(float64)),
				int(arr[2].(float64)),
				int(arr[3].(float64)),
				int(arr[4].(float64)),
			)
		}
	}

	return graph, nil
}
