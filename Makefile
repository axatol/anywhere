all: clean server

server: clean
	go mod tidy
	go build -o ./bin/server ./cmd/server/main.go

docker: server
	docker build -t anywhere .

clean:
	rm -rf ./bin
