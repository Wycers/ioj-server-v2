package ws

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/infinity-oj/server-v2/internal/pkg/utils/json"

	"github.com/infinity-oj/server-v2/internal/pkg/eventBus"
	"github.com/infinity-oj/server-v2/pkg/models"

	"github.com/gorilla/websocket"
)

type Actuator struct {
	hub *Hub
	*Client
}

type User struct {
	hub *Hub
	*Client
}

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 8) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin: func(r *http.Request) bool {
	// 	return true
	// },
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *User) readPump() {
	defer func() {
		c.conn.Close()
		c.hub.unregister <- c
	}()
	fmt.Println("gg")
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPingHandler(func(string) error {
		fmt.Println("pinged")
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	c.conn.SetPongHandler(func(string) error {
		fmt.Println("ponged")
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))

		fmt.Println(string(message))
		//c.hub.broadcast <- message
	}
	fmt.Println("close")
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *User) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer c.conn.Close()

	bus := eventBus.New()

	newTaskChan := make(chan *models.Task)

	if err := bus.Subscribe("task:new", func(task *models.Task) {
		fmt.Println("task id", task.TaskId)
		newTaskChan <- task
	}); err != nil {
		fmt.Println(err)
		return
	}

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case task := <-newTaskChan:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			//if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
			//	return
			//}
			if data, err := json.Marshal(task); err != nil {
				fmt.Println("show task with error", err)
				return
			} else {
				if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
					fmt.Println("show task with error", err)
					return
				} else {
					fmt.Println("show task")
				}
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				fmt.Println("ping with error", err)
				return
			}
		}
	}
}

// ServeW handles websocket requests from the peer.
func (hub *Hub) ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := &Client{
		conn: conn,
		send: make(chan []byte, 256),
	}
	user := &User{
		hub:    hub,
		Client: client,
	}

	user.hub.register <- user

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go user.writePump()
	go user.readPump()
}
