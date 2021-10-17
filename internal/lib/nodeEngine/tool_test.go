package nodeengine

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/infinity-oj/server-v2/internal/lib/nodeengine/scene"
)

func TestScene(t *testing.T) {
	body, err := os.ReadFile("scene/blocks.json")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	blocksJSONStr := string(body)

	blocks := scene.NewBlocksDefinition(blocksJSONStr)

	body, err = os.ReadFile("scene/scene.json")
	if err != nil {
		log.Fatalf("unable to read file: %v", err)
	}
	sceneJSONStr := string(body)

	scene := scene.NewScene(sceneJSONStr)

	graph, err := NewGraphByScene(blocks, scene)
	fmt.Println(graph)

	for _, v := range graph.Run() {
		fmt.Println(v.Type)
	}
}
