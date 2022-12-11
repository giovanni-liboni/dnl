# Distributed Notification Library

DNL is a library for sending notifications to multiple instances of the same service or beween different services. It is 
designed to be used in a distributed environment, where the services are running on different machines or pods (k8s).
The main goal of this library is to provide a simple and easy to use interface.

## Installation

```bash
go get github.com/alexandrevicenzi/dnl
```


## Usage

The basic idea is to provide a backend provider, that manages the distribution of the notifications to the different
instances of the notification service.


## Example

```go
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

	// Add a channel to the DWL
	ch1.AddChannel("channel1")
	ch2.AddChannel("channel2")

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
```