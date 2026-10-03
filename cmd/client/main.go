package main

import (
	"fmt"
	"os"
	"os/signal"

	gamelogic "github.com/absoluteKeven/go-pub-sub/internal/gamelogic"
	pubsub "github.com/absoluteKeven/go-pub-sub/internal/pubsub"
	routing "github.com/absoluteKeven/go-pub-sub/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	connS := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connS)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer conn.Close()

	name, err := gamelogic.ClientWelcome()

	_, _, errps := pubsub.DeclareAndBind(conn, "peril_direct", (routing.PauseKey + "." + name), routing.PauseKey, pubsub.Transient)
	if errps != nil {
		fmt.Printf("%s\n", err)
		return
	}

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("Shutting gracefully.")
}
