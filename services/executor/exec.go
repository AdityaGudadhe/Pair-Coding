package main

import (
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func populateMessage(m []byte) (*Message, error) {
	var message Message

	err := json.Unmarshal(m, message)
	if err != nil {
		return nil, err
	}

	return &message, nil
}

func operator(message *amqp.Delivery) {
	m, err := populateMessage(message.Body)
	if err != nil {
		log.Printf("Error while parsing message: %s", err)
		return
	}
	fmt.Print(m.problemId)
}

func consumeMsg(conn *amqp.Connection, consumerID int) {
	msgs, ch, err := openHelloChannel(conn, consumerID)
	if err != nil {
		return
	}
	defer ch.Close()
	for msg := range msgs {
		operator(&msg)
	}
	log.Printf("consumer %d: stopped", consumerID)
}
