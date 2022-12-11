package main

import (
	"dnl"
	"fmt"
	"log"
)

// Simple example of how to use the DWL

func main() {
	// Create a new DWL
	provider1 := dnl.NewProviderRedis("localhost:6379")
	ch1 := dnl.NewWithProvider(provider1)

	// Create a new channel
	provider2 := dnl.NewProviderRedis("localhost:6379")
	ch2 := dnl.NewWithProvider(provider2)

	// Create a new channel
	provider3 := dnl.NewProviderRedis("localhost:6379")
	ch3 := dnl.NewWithProvider(provider3)

	// Add a channel to the DWL
	ch1.AddChannel("channel1")
	ch2.AddChannel("channel2")
	ch3.AddChannel("channel3")

	// Listen to the channel
	go listenToChannel(ch1, "channel1")
	go listenToChannel(ch2, "channel2")

	// Send a message to the channel
	err := ch1.Send("channel1", "Hello World!")
	if err != nil {
		log.Println(err)

		return
	}

	err = ch2.Send("channel1", "Hello World from another channel!")
	if err != nil {
		log.Println(err)

		return
	}

	err = ch2.Send("channel1", "Hello World from another goroutine passing from Redis!")
	if err != nil {
		log.Println(err)

		return
	}

	// Wait forever
	select {}
}

func listenToChannel(channels dnl.DNL, id string) {
	for {
		select {
		case msg := <-channels.GetChannel(id):
			fmt.Println("Received message for channel", id, ":", msg)
		}
	}
}
