## go-chat

A simple golang chat server and client.

## Usage

Start a chat server:

```
cd cmd/server
go run main.go
```

In a new tab/window, start a chat client:

```
cd cmd/client
go run main.go
```

## Build for Windows

```
GOOS=windows GOARCH=amd64 go build
```

## To Do
- Fix Windows new lines
- UI - https://github.com/marcusolsson/tui-go
