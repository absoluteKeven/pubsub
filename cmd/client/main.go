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

	go func() {
		<-c
		fmt.Println("Shutting gracefully.")
		conn.Close()
		os.Exit(0)
	}()

	gameState := gamelogic.NewGameState(name)

cmdloop:
	for {
		cmd := gamelogic.GetInput()

		switch cmd[0] {
		case "spawn":
			err := gameState.CommandSpawn(cmd[1:])
			if err != nil {
				fmt.Printf("Error spawning unit: %s\n", err)
			}
		case "move":
			move, err := gameState.CommandMove(cmd[1:])
			if err != nil {
				fmt.Printf("Error moving unit: %s\n", err)
				break
			}
			fmt.Printf("Moved %i to %s\n", move.Units[0].ID, move.ToLocation)
		case "status":
			gameState.CommandStatus()
			break
		case "help":
			gamelogic.PrintClientHelp()
			break
		case "spam":
			fmt.Println("Spamming not alowed.")
			break
		case "quit":
			fmt.Println("Exiting.")
			break cmdloop
		default:
			fmt.Println("Unknown command.")
		}
	}
}
