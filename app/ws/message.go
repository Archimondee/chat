package ws

import (
	"chat/app/models/entity"
	"encoding/json"
	"log"
)

const UserJoinedAction = "user-action"
const UserLeftAction = "user-left"
const UserOnline = "user-online"
const SendMessage = "send-message"

type Message struct {
	Action    string `json:"action"`
	Message   string `json:"message"`
	Recipient string `json:"recipient"`
	//Target  *Room       `json:"target"`
	Sender string `json:"sender"`
}

type ListOnlineMessage struct {
	Action string         `json:"action"`
	Users  []*entity.User `json:"users"`
}

func (message *Message) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}

func (message *ListOnlineMessage) encode() []byte {
	json, err := json.Marshal(message)
	if err != nil {
		log.Println(err)
	}

	return json
}

//func (message *Message) UnmarshalJSON(data []byte) error {
//	type Alias Message
//	msg := &struct {
//		Sender Client `json:"sender"`
//		*Alias
//	}{
//		Alias: (*Alias)(message),
//	}
//	if err := json.Unmarshal(data, &msg); err != nil {
//		return err
//	}
//	message.Sender = &msg.Sender
//	return nil
//}
