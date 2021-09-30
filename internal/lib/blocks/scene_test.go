package blocks

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("New", func(t *testing.T) {

		body, err := os.ReadFile("scene.json")
		if err != nil {
			log.Fatalf("unable to read file: %v", err)
		}
		sceneJSONStr := string(body)

		scene := New(sceneJSONStr)
		fmt.Println(scene)

		//fmt.Println(string(body))
	})
}
