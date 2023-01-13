run:
	go run .

build:
	go build .

build-linux:
	env GOOS=linux GOARCH=amd64 go build -o odop .
