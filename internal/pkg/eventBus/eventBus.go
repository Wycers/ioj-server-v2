package eventBus

import (
	"fmt"
	"sync"

	"github.com/asaskevich/EventBus"
)

var once sync.Once

type Bus = EventBus.Bus

var bus Bus

func New() Bus {
	once.Do(func() {
		fmt.Println("???????")
		bus = EventBus.New()
	})
	fmt.Println(&bus)
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
