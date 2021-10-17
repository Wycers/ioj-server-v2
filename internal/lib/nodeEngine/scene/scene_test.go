package scene

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("NewScene", func(t *testing.T) {

		body, err := os.ReadFile("scene.json")
		if err != nil {
			log.Fatalf("unable to read file: %v", err)
		}
		sceneJSONStr := string(body)

		scene := NewScene(sceneJSONStr)
		fmt.Println(scene)
		//fmt.Println(string(body))
	})
}
