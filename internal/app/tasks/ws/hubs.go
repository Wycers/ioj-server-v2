package ws

import (
	"fmt"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered actuators.
	actuators map[*Actuator]bool

	taskTopics map[string]*unicastSubject

	// Register requests from the actuators.
	register chan *Actuator

	// Unregister requests from actuators.
	unregister chan *Actuator
}

var once sync.Once
var hub *Hub

func NewHub() *Hub {

	once.Do(func() {
		hub = &Hub{
			actuators:  make(map[*Actuator]bool),
			taskTopics: make(map[string]*unicastSubject),
			register:   make(chan *Actuator),
			unregister: make(chan *Actuator),
		}
	})
	go hub.run()
	return hub
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			fmt.Println("??")
			hub.actuators[client] = true
		case client := <-hub.unregister:
			if _, ok := hub.actuators[client]; ok {
				delete(hub.actuators, client)
				close(client.send)
			}
		}
	}
}
