# Distributed Notification Library

DNL is a library designed for sending notifications to multiple instances of the same service or between different services. It's particularly useful in a distributed environment where services run on separate machines or pods (k8s). The primary aim of this library is to offer a simple and easy-to-use interface.

## Installation

```bash
go get github.com/giovanni-liboni/dnl
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

// Simple example of how to use the DNL

func main() {
	// Create a new DWL
	provider1 := dnl.NewProviderRedis("localhost:6379")
	ch1 := dnl.NewWithProvider(provider1)

	// Create a new channel
	provider2 := dnl.NewProviderRedis("localhost:6379")
	ch2 := dnl.NewWithProvider(provider2)

	// Add a channel to the DNL
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

Output:
```
Received message for channel channel1 : Hello World!
Received message for channel channel1 : Hello World from another channel!
Received message for channel channel1 : Hello World from another goroutine passing from Redis!
```

## License

The project is licensed under the MIT license. See the [LICENSE](LICENSE) file for more details.

## Contributing

Feel free to contribute to this project. Any help is appreciated.
