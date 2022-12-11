# Server-sent events example

In this example we will create a simple server-sent events application using the DNL. The application will have 
two replicas. Each replica can send a message to the other replica if the client is connected to that replica.

## Running the example

To run the example, first run the docker-compose file in the repository:

```bash
$ docker-compose up
```

Then, in a separate terminal, run the following command:

```bash
$ make run-server-1
```

this will run a server on port 3001. To run the second server, run the following command:

```bash
$ make run-server-2
```

this will run a server on port 3002.

## Testing the example

Using a Websocket client (e.g. websocat), connect to the server on `ws://localhost:3002/ws/1`

```bash
$ websocat ws://127.0.0.1:3002/ws/1
```

Then, in a separate terminal, send a body to the server on port 3001:

```bash
$ make send-message-1
```

You should see the message in the Websocket client. As we're sending the message to the other replica, the message is
using the library to transfer the message from one server to the other using the provider. 