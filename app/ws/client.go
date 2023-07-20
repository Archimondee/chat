package ws

import (
	"chat/app/interfaces"
	"chat/app/models/entity"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	// Max wait time when writing message to peer
	writeWait = 10 * time.Second
	// Max time till next pong from peer
	pongWait = 60 * time.Second
	// Send ping interval, must be less then pong wait time
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 10000
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	conn   *websocket.Conn
	server *Server
	send   chan []byte
	user   *entity.User
}

func ServeWebsocket(server *Server, w http.ResponseWriter, r *http.Request, userRepository interfaces.UserRepository) {
	uuid, ok := r.URL.Query()["chat"]
	if !ok || len(uuid[0]) < 1 {
		log.Println("Url Param 'chat' is missing")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	res, err := userRepository.FindUserByUuid(uuid[0])
	if err != nil {
		fmt.Println("Error", err)
	}

	client := newClient(conn, server, res)

	go client.writePump()
	go client.readPump()

	server.register <- client
}

func newClient(conn *websocket.Conn, server *Server, user *entity.User) *Client {
	return &Client{
		server: server,
		conn:   conn,
		user:   user,
		send:   make(chan []byte, 256),
	}
}

func (client *Client) readPump() {
	defer func() {
		client.disconnect()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	// Start endless read loop, waiting for messages from client
	for {
		_, jsonMessage, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("unexpected close error: %v", err)
			}
			break
		}

		client.handleNewMessage(jsonMessage)
	}

}

func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The WsServer closed the channel.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Attach queued chat messages to the current websocket message.
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (client *Client) handleNewMessage(jsonMessage []byte) {

	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}

	message.Sender = *client.user
	fmt.Println("dataa", client)

	//switch message.Action {
	//case SendMessageAction:
	//	roomID := message.Target.GetId()
	//	if room := client.wsServer.findRoomByID(roomID); room != nil {
	//		room.broadcast <- &message
	//	}
	//
	//case JoinRoomAction:
	//	client.handleJoinRoomMessage(message)
	//
	//case LeaveRoomAction:
	//	client.handleLeaveRoomMessage(message)
	//
	//case JoinRoomPrivateAction:
	//	client.handleJoinRoomPrivateMessage(message)
	//}

}

func (client *Client) disconnect() {
	client.server.unregister <- client
	close(client.send)
	client.conn.Close()
}