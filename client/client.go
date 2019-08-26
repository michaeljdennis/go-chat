package client

import (
	"bufio"
	"log"
	"net"
	"os"

	"github.com/michaeljdennis/go-chat/client/user"
	"github.com/michaeljdennis/go-chat/message"
)

// Client ...
type Client struct{}

// New ...
func New() *Client {
	return &Client{}
}

// Join ...
func (c *Client) Join() {
	var err error

	// Dial up to AOL
	conn, err := net.Dial("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()

	// Create new user
	u := user.New()

	// Ask for username
	stdout := bufio.NewWriter(os.Stdout)
	_, err = message.WriteMessage(stdout, "Enter username:")
	if err != nil {
		log.Fatalln(err)
	}

	// Get username
	stdin := bufio.NewReader(os.Stdin)
	u.Name, err = message.ReadMessage(stdin)
	if err != nil {
		log.Fatalln(err)
	}

	// Send username
	writer := bufio.NewWriter(conn)
	_, err = message.WriteMessage(writer, u.Name)
	if err != nil {
		log.Fatalln(err)
	}

	// Receive UUID response
	reader := bufio.NewReader(conn)
	u.ID, err = message.ReadMessage(reader)
	if err != nil {
		log.Fatalln(err)
	}

	// Listen for messages from broadcast
	go func() {
		for {
			reader := bufio.NewReader(conn)
			msg, err := message.ReadMessage(reader)
			if err != nil {
				log.Fatalln(err)
			}

			writer := bufio.NewWriter(os.Stdout)
			message.WriteMessage(writer, msg)
		}
	}()

	for {
		// Get message from STDIN
		stdin := bufio.NewReader(os.Stdin)
		msg, err := message.ReadMessage(stdin)
		if err != nil {
			log.Fatalln(err)
		}

		// Send message
		writer := bufio.NewWriter(conn)
		_, err = message.WriteMessage(writer, msg)
		if err != nil {
			log.Fatalln(err)
		}
	}
}
