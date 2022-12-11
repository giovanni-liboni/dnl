run-server-1:
	PORT=3001 go run main.go

run-server-2:
	PORT=3002 go run main.go

send-message-1:
	curl -X POST -d "message=hello" http://localhost:3001/notify/1
