package redis

import (
	"context"
	"fmt"
	"testing"

	"go.uber.org/zap"

	"github.com/infinity-oj/server-v2/internal/pkg/configs"
)

func TestNew(t *testing.T) {
	v, _ := configs.New("configs/server.yaml")
	fmt.Println(v.Get("redis"))
	opts, err := NewOptions(v, zap.NewNop())
	fmt.Println(opts, err)
	r := New(opts)
	r.Set(context.Background(), "y", 1, 0)
	fmt.Println(r.Get(context.Background(), "y").Uint64())
}
