package ws

import (
	"fmt"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered users.
	users map[*User]bool

	// Register requests from the actuators.
	register chan *User

	// Unregister requests from actuators.
	unregister chan *User
}

var once sync.Once
var hub *Hub

func NewHub() *Hub {
	once.Do(func() {
		hub = &Hub{
			users:      make(map[*User]bool),
			register:   make(chan *User),
			unregister: make(chan *User),
		}
	})
	go hub.run()
	return hub
}

func (hub *Hub) run() {
	for {
		select {
		case client := <-hub.register:
			fmt.Println("register", client)
			hub.users[client] = true
		case client := <-hub.unregister:
			fmt.Println("unregister", client)
			if _, ok := hub.users[client]; ok {
				delete(hub.users, client)
				close(client.send)
			}
		}
	}
}
