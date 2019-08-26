package main

import "github.com/michaeljdennis/go-chat/room"

func main() {
	// Create chat room
	r := room.New()
	r.Open()
}
