package utils

import (
	"github.com/streadway/amqp"
	"log"
)

// var AmqpConn *amqp.Connection
var AmqpChannel *amqp.Channel
var AmqpQueue *amqp.Queue

func ConnectAmqp(amqpUrl string, amqpQueue string, routing string, exchange string) {
	conn, err := amqp.Dial(amqpUrl)
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ:", err)
	}
	//defer AmqpConn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel:", err)
	}
	//defer AmqpChannel.Close()

	q, err := ch.QueueDeclare(
		amqpQueue, // queue name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare a queue:", err)
	}
	//
	err = ch.ExchangeDeclare(
		exchange, // name
		"direct", // type
		false,    // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatal("Failed to declare an exchange:", err)
	}

	err = ch.QueueBind(
		amqpQueue, //queue name
		routing,   //routing
		exchange,  //exchange
		false,
		nil,
	)

	if err != nil {
		log.Fatal("Failed to bind the queue:", err)
	}

	AmqpQueue = &q
	AmqpChannel = ch

}
