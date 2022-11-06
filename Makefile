all: clean server

server: clean
	go mod tidy
	go build -o ./bin/server ./cmd/server/main.go

clean:
	rm -rf ./bin
