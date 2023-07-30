package ws

import (
	"chat/app/interfaces"
	"chat/config"
	"chat/utils"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/streadway/amqp"
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
	conn              *websocket.Conn
	server            *Server
	send              chan []byte
	uuid              string
	messageRepository interfaces.MessageRepository
	//user   *entity.User

}

func ServeWebsocket(server *Server, w http.ResponseWriter, r *http.Request, messageRepository interfaces.MessageRepository) {
	id, ok := r.URL.Query()["chat"]
	if !ok || len(id[0]) < 1 {
		log.Println("Url Param 'chat' is missing")
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	//res, err := userRepository.FindUserByUuid(uuid[0])
	//if err != nil {
	//	fmt.Println("Error", err)
	//}

	client := newClient(conn, server, id[0], messageRepository)

	go client.writePump()
	go client.readPump()
	go client.consumeMessage()
	go server.RunRoomRepository(client)

	server.register <- client
}

func newClient(conn *websocket.Conn, server *Server, uuid string, messageRepository interfaces.MessageRepository) *Client {
	return &Client{
		server:            server,
		conn:              conn,
		uuid:              uuid,
		send:              make(chan []byte, 256),
		messageRepository: messageRepository,
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
	var message interfaces.Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		log.Printf("Error on unmarshal JSON message %s", err)
		return
	}

	//message.Sender = client.uuid

	switch message.Action {
	case SendMessage:
		client.sendMessage(&message.Sender, message.Recipient, []byte(message.Message))
	case SendGroupMessage:
		roomClient, ok := client.server.rooms[message.RoomId]
		if ok {
			roomClient.broadcast <- &message
		}
	}

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

func (client *Client) sendMessage(sender *string, id string, message []byte) {
	//recipientClient, ok := client.server.clients[id]
	loadedConfig, _ := config.LoadConfig(".")

	m := &interfaces.Message{
		Action:    SendMessage,
		Sender:    *sender,
		Recipient: id,
		Message:   string(message),
		Uuid:      uuid.New(),
		Status:    "not_sent",
	}
	client.messageRepository.CreateMessage(*m)
	jsonByte, err := json.Marshal(m)
	if err != nil {
		log.Println(err)
	}

	err = utils.AmqpChannel.Publish(
		loadedConfig.AmqpExchange, //exchange
		loadedConfig.AmqpRouting,  // routing key (queue name)
		false,                     // mandatory
		false,                     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        jsonByte,
		},
	)
	if err != nil {
		fmt.Println("Failed to publish messages : ", err)
	}
}

func (client *Client) consumeMessage() {
	loadedConfig, err := config.LoadConfig("../..")
	if err != nil {
		log.Println(err)
	}
	msgs, err := utils.AmqpChannel.Consume(
		loadedConfig.AmqpQueue, // queue
		"",                     // consumer
		false,                  // auto-ack
		false,                  // exclusive
		false,                  // no-local
		false,                  // no-wait
		nil,
	)

	if err != nil {
		log.Fatal("Failed to register a consumer:", err)
	}
	go func() {
		for msg := range msgs {
			BroadcastClient(msg, client)
		}
	}()

}

func BroadcastClient(msg amqp.Delivery, client *Client) {
	var body *interfaces.Message
	if err := json.Unmarshal(msg.Body, &body); err != nil {
		log.Println(err)
	}
	recipientClient, ok := client.server.clients[body.Recipient]
	if ok {
		body.Status = "sent"
		errUpdate := client.messageRepository.UpdateMessage(*body)
		if errUpdate != nil {
			log.Println(errUpdate)
		}

		client.server.clients[recipientClient.uuid].send <- Encode(body)

		err := msg.Ack(false)
		if err != nil {
			log.Println(err)
		}

	} else {
		err := msg.Nack(false, true)
		if err != nil {
			log.Println(err)
		}
	}
}
