package main

import (
	"bufio"
	"dnl"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/websocket/v2"
)

var index = []byte(`<!DOCTYPE html>
<html>
<body>
<h1>SSE Messages</h1>
<div id="result"></div>
<script>
if(typeof(EventSource) !== "undefined") {
  var source = new EventSource("http://127.0.0.1:3000/sse");
  source.onmessage = function(event) {
    document.getElementById("result").innerHTML += event.data + "<br>";
  };
} else {
  document.getElementById("result").innerHTML = "Sorry, your browser does not support server-sent events...";
}
</script>
</body>
</html>
`)

func main() {
	app := fiber.New()

	// CORS for external resources
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowHeaders:     "Cache-Control",
		AllowCredentials: true,
	}))

	// Initialize the NDL
	provider := dnl.NewProviderRedis("localhost:6379")
	channels := dnl.NewWithProvider(provider)

	app.Use("/ws", func(c *fiber.Ctx) error {
		// IsWebSocketUpgrade returns true if the client
		// requested upgrade to the WebSocket protocol.
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})

	app.Post("/notify/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		log.Println("Notify", id)
		log.Println("Body", string(c.Body()))
		err := channels.Send(id, string(c.Body()))
		if err != nil {
			return err
		}

		return c.SendStatus(fiber.StatusOK)
	})

	app.Get("/ws/:id", websocket.New(func(c *websocket.Conn) {
		// c.Locals is added to the *websocket.Conn
		log.Println(c.Locals("allowed"))  // true
		log.Println(c.Params("id"))       // 123
		log.Println(c.Query("v"))         // 1.0
		log.Println(c.Cookies("session")) // ""
		id := c.Params("id")

		// Add the channel to the NDL
		channels.AddChannel(id)
		defer channels.RemoveChannel(id)

		// websocket.Conn bindings https://pkg.go.dev/github.com/fasthttp/websocket?tab=doc#pkg-index
		var (
			msg string
			err error
		)
		log.Println("Listening for messages on channel", id)
		for {
			select {
			case msg = <-channels.GetChannel(id):
				// Write message to connection
				if err = c.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
					log.Println("write:", err)
					return
				}
			}
		}
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		c.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)

		return c.Status(fiber.StatusOK).Send(index)
	})

	app.Get("/sse/:id", func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		id := c.Params("id")
		// Add the channel to the NDL
		channels.AddChannel(id)
		defer channels.RemoveChannel(id)

		c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
			fmt.Println("WRITER")
			var i int
			for {
				i++
				msg := fmt.Sprintf("%d - the time is %v", i, time.Now())
				fmt.Fprintf(w, "data: Message: %s\n\n", msg)
				fmt.Println(msg)

				fmt.Println("ID:", id)

				err := w.Flush()
				if err != nil {
					// Refreshing page in web browser will establish a new
					// SSE connection, but only (the last) one is alive, so
					// dead connections must be closed here.
					fmt.Printf("Error while flushing: %v. Closing http connection.\n", err)

					break
				}
				time.Sleep(2 * time.Second)
			}
		})

		return nil
	})

	// Get the port from the environment
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Start the server
	log.Fatal(app.Listen(":" + port))
	// Access the websocket server: ws://localhost:3000/ws/123?v=1.0
	// https://www.websocket.org/echo.html
}
