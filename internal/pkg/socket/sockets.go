package socket

import (
	socketio "github.com/googollee/go-socket.io"
)

func Handler() {
	server := socketio.NewServer(nil)

	return server
}
