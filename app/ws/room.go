package ws

import (
	"chat/app/interfaces"
	"fmt"
)

type Room struct {
	uuid       string `json:"uuid"`
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan *interfaces.Message
}

func NewRoom(uuid string) *Room {
	return &Room{
		uuid:       uuid,
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *interfaces.Message),
	}
}

func (room *Room) RunRoom() {
	for {
		select {

		case client := <-room.register:
			room.registerClientInRoom(client)

		case client := <-room.unregister:
			room.unregisterClientInRoom(client)

		case message := <-room.broadcast:
			fmt.Println("data 1234", message)

			room.broadcastToClientsInRoom(message)
		}
	}
}

func (room *Room) registerClientInRoom(client *Client) {
	room.clients[client.uuid] = client
}

func (room *Room) unregisterClientInRoom(client *Client) {
	if _, ok := room.clients[client.uuid]; ok {
		delete(room.clients, client.uuid)
	}
}

func (room *Room) broadcastToClientsInRoom(message *interfaces.Message) {
	fmt.Println("data", message)
	//Will add to rabbit mq
	//for _, client := range room.clients {
	//
	//	client.send <- message
	//}
}
