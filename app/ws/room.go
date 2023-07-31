package ws

import (
	"chat/app/interfaces"
	"chat/utils"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"log"
)

type Room struct {
	uuid              string `json:"uuid"`
	clients           map[string]*Client
	register          chan *Client
	unregister        chan *Client
	broadcast         chan *interfaces.Message
	roomRepository    interfaces.RoomRepository
	messageRepository interfaces.MessageRepository
}

func NewRoom(uuid string, roomRepository interfaces.RoomRepository, messageRepository interfaces.MessageRepository) *Room {
	return &Room{
		uuid:              uuid,
		clients:           make(map[string]*Client),
		register:          make(chan *Client),
		unregister:        make(chan *Client),
		broadcast:         make(chan *interfaces.Message),
		roomRepository:    roomRepository,
		messageRepository: messageRepository,
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

	resRoom, err := room.roomRepository.FindRoomById(message.RoomId)
	if err != nil {
		fmt.Println("Error ", err)
	}
	m := &interfaces.Message{
		Action: SendGroupMessage,
		Sender: message.Sender,
		//Recipient: participant.UserId.String(),
		Message: message.Message,
		Uuid:    uuid.New(),
		Status:  "sent",
		RoomId:  message.RoomId,
	}

	room.messageRepository.CreateMessage(*m)

	for _, participant := range resRoom.Participants {

		if participant.UserId.String() != message.Sender {

			m.Recipient = participant.UserId.String()
			jsonByte, err := json.Marshal(m)
			if err != nil {
				log.Println(err)
			}
			err = utils.AmqpGroupChannel.Publish(
				"message_group_exchange",
				"message_group_routing", // routing key (queue name)
				false,                   // mandatory
				false,                   // immediate
				amqp.Publishing{
					ContentType: "text/plain",
					Body:        jsonByte,
				},
			)
			if err != nil {
				fmt.Println("Failed to publish messages : ", err)
			}

		}

		//cl, ok := room.clients[participant.UserId.String()]
		//if ok {
		//	message.Status = "sent"
		//}
		//errUpdate := room.messageRepository.UpdateMessage(*message)
		//if errUpdate != nil {
		//	log.Println(errUpdate)
		//}
	}

	//Will add to rabbit mq
	//for _, client := range room.clients {
	//
	//	client.send <- message
	//}
}

//func (room *Room) consumeGroupMessage() {
//	loadedConfig, err := config.LoadConfig("../..")
//	if err != nil {
//		log.Println(err)
//	}
//	msgs, err := utils.AmqpChannel.Consume(
//		loadedConfig.AmqpGroupQueue, // queue
//		"",                          // consumer
//		false,                       // auto-ack
//		false,                       // exclusive
//		false,                       // no-local
//		false,                       // no-wait
//		nil,
//	)
//
//	if err != nil {
//		log.Fatal("Failed to register a consumer:", err)
//	}
//	go func() {
//		for msg := range msgs {
//			BroadcastClientGroup(msg, room)
//		}
//	}()
//
//}

//func BroadcastClientGroup(msg amqp.Delivery, room *Room) {
//	var body *interfaces.Message
//	if err := json.Unmarshal(msg.Body, &body); err != nil {
//		log.Println(err)
//	}
//	recipientClient, ok := room.clients[body.Sender]
//	if ok {
//		body.Status = "sent"
//		errUpdate := room.messageRepository.UpdateMessage(*body)
//		if errUpdate != nil {
//			log.Println(errUpdate)
//		}
//
//		room.clients[recipientClient.uuid].send <- Encode(body)
//
//		err := msg.Ack(false)
//		if err != nil {
//			log.Println(err)
//		}
//
//	} else {
//		err := msg.Nack(false, true)
//		if err != nil {
//			log.Println(err)
//		}
//	}
//}
