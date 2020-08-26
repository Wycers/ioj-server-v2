package nodeEngine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
				fmt.Println("==> value not string:", v)
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
