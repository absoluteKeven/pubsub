package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/absoluteKeven/go-pub-sub/internal/gamelogic"
	pubsub "github.com/absoluteKeven/go-pub-sub/internal/pubsub"
	routing "github.com/absoluteKeven/go-pub-sub/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	connS := "amqp://guest:guest@localhost:5672/"

	conn, err := amqp.Dial(connS)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}

	fmt.Println("Connection Successful...")

	gamelogic.PrintServerHelp()

	for {
		cmd := gamelogic.GetInput()

		switch cmd[0] {
		case routing.PauseKey:
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
		case "resume":
			pubsub.PublishJSON(ch, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
		case "quit":
			fmt.Println("Exiting.")
			break
		default:
			fmt.Println("Unknown command.")
		}
	}

	c := make(chan os.Signal, 1)

	signal.Notify(c, os.Interrupt)
	<-c
	fmt.Println("Shutting gracefully.")
}
