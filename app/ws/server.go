package ws

import (
	"chat/app/interfaces"
	"chat/app/models/entity"
	"fmt"
)

type Server struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte

	userRepository interfaces.UserRepository
	users          []*entity.User
}

func NewWebsocketServer(userRepository interfaces.UserRepository) *Server {
	s := &Server{
		clients:        make(map[*Client]bool),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan []byte),
		userRepository: userRepository,
	}
	//res, err := userRepository.FindUserAll()
	//if err != nil {
	//	fmt.Println("Error finduserall", err)
	//}

	//s.users = res

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
			server.broadcastToClients(message)
		}
	}
}

func (server *Server) registerClient(client *Client) {
	server.userRepository.UpdateStatus("online", client.user.Uuid.String())
	server.listOnlineClients(client)
	server.clients[client] = true
	server.users = append(server.users, client.user)
}

func (server *Server) unregisterClient(client *Client) {
	if _, ok := server.clients[client]; ok {
		delete(server.clients, client)
		for i, user := range server.users {
			if user.Uuid == client.user.Uuid {
				server.userRepository.UpdateStatus("offline", client.user.Uuid.String())
				server.listOnlineClients(client)

				server.users[i] = server.users[len(server.users)-1]
				server.users = server.users[:len(server.users)-1]
				break
			}
		}

	}
}

func (server *Server) notifyClientJoined(client *Client) {
	message := &Message{
		Action: UserJoinedAction,
		Sender: *client.user,
	}

	server.broadcastToClients(message.encode())
}

func (server *Server) broadcastToClients(message []byte) {
	for client := range server.clients {
		client.send <- message
	}
}

func (server *Server) notifyClientLeft(client *Client) {
	message := &Message{
		Action: UserLeftAction,
		Sender: *client.user,
	}

	server.broadcastToClients(message.encode())
}

func (server *Server) listOnlineClients(client *Client) {
	res, err := server.userRepository.FindUserAll()
	if err != nil {
		fmt.Println("Error finduserall", err)
	}

	message := &ListOnlineMessage{
		Action: UserOnline,
		Users:  res,
	}
	server.broadcastToClients(message.encode())
}
