package ws

import (
	"chat/app/interfaces"
	"fmt"
)

type Server struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte

	rooms map[string]*Room

	userRepository    interfaces.UserRepository
	users             []*string
	roomRepository    interfaces.RoomRepository
	messageRepository interfaces.MessageRepository
}

func NewWebsocketServer(userRepository interfaces.UserRepository, roomRepository interfaces.RoomRepository, messageRepository interfaces.MessageRepository) *Server {
	s := &Server{
		clients:           make(map[string]*Client),
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		broadcast:         make(chan []byte),
		userRepository:    userRepository,
		rooms:             make(map[string]*Room),
		roomRepository:    roomRepository,
		messageRepository: messageRepository,
	}

	return s
}

func (server *Server) Run() {
	for {
		select {
		case client := <-server.register:
			server.registerClient(client)
		case client := <-server.unregister:
			server.unregisterClient(client)
		case message := <-server.broadcast:
			server.broadcastToAllClientsConnected(message)
		}
	}
}

func (server *Server) registerClient(client *Client) {
	server.userRepository.UpdateStatus("online", client.uuid)
	server.listOnlineClients(client)
	server.clients[client.uuid] = client
	server.users = append(server.users, &client.uuid)
}

func (server *Server) unregisterClient(client *Client) {
	if _, ok := server.clients[client.uuid]; ok {
		delete(server.clients, client.uuid)
		for i, user := range server.users {
			if user == &client.uuid {
				server.userRepository.UpdateStatus("offline", client.uuid)
				server.listOnlineClients(client)

				server.users[i] = server.users[len(server.users)-1]
				server.users = server.users[:len(server.users)-1]
				break
			}
		}
	}
}

//func (server *Server) notifyClientJoined(client *Client) {
//	message := &Message{
//		Action: UserJoinedAction,
//		Sender: *client.user,
//	}
//
//	server.broadcastToAllClientsConnected(message.encode())
//}

func (server *Server) broadcastToAllClientsConnected(message []byte) {
	user, err := server.userRepository.FindUserAll()
	if err != nil {
		fmt.Println("Error")
	}

	for _, i := range user {
		recipient, ok := server.clients[i.Uuid.String()]
		if ok {
			server.clients[recipient.uuid].send <- message
		}
	}
}

//func (server *Server) notifyClientLeft(client *Client) {
//	message := &Message{
//		Action: UserLeftAction,
//		Sender: *client.user,
//	}
//
//	server.broadcastToAllClientsConnected(message.encode())
//}

func (server *Server) listOnlineClients(client *Client) {
	res, err := server.userRepository.FindUserAll()
	if err != nil {
		fmt.Println("Error finduserall", err)
	}

	message := &ListOnlineMessage{
		Action: UserOnline,
		Users:  res,
	}
	server.broadcastToAllClientsConnected(message.encode())
}

func (server *Server) RunRoomRepository(client *Client) {

	rooms, err := server.roomRepository.GetAllRoom()
	if err != nil {
		fmt.Println("Error ", err)
	}

	for _, room := range rooms {
		r := NewRoom(room.Uuid.String(), server.roomRepository, server.messageRepository)
		go r.RunRoom()
		//go r.consumeGroupMessage()
		server.rooms[room.Uuid.String()] = r
		server.rooms[room.Uuid.String()].register <- client

	}
}
