package eventBuses

import (
	"fmt"
	"github.com/asaskevich/EventBus"
	"sync"
)

var once sync.Once

var bus EventBus.Bus

func New() EventBus.Bus {
	once.Do(func() {
		bus = EventBus.New()
	})
	return bus
}

func calculator(a int, b int) {
	fmt.Printf("%d\n", a+b)
}

func XD() {
	bus.Subscribe("main:calculator", calculator)
	bus.Publish("main:calculator", 20, 40)
	bus.Unsubscribe("main:calculator", calculator)
}
