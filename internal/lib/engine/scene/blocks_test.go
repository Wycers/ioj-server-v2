package scene

import (
	"fmt"
	"log"
	"os"
	"testing"
)

func TestNewBlocksDefinition(t *testing.T) {
	t.Run("NewScene", func(t *testing.T) {

		body, err := os.ReadFile("blocks.json")
		if err != nil {
			log.Fatalf("unable to read file: %v", err)
		}
		blocksJSONStr := string(body)

		blocks := NewBlocksDefinition(blocksJSONStr)
		fmt.Println(blocks)
		//fmt.Println(string(body))
	})
}
