package main

import "github.com/michaeljdennis/go-chat/client"

func main() {
	// Create client
	c := client.New()
	c.Join()
}
