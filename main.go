package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
)

func main() {
	chatRoom := NewChatRoom()
	chatRoom.Listen()

	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln(err)
		}
		chatRoom.Join(conn)
	}
}

// ChatRoom ...
type ChatRoom struct {
	clients []*Client
	// messages from a client to be sent out
	hotline chan ClientMessage
}

// NewChatRoom ...
func NewChatRoom() *ChatRoom {
	chatRoom := &ChatRoom{}
	chatRoom.hotline = make(chan ClientMessage)

	return chatRoom
}

// Join ...
func (chatRoom *ChatRoom) Join(conn net.Conn) {
	id := len(chatRoom.clients) + 1
	client := NewClient(conn, id)
	chatRoom.clients = append(chatRoom.clients, client)

	go func() {
		for {
			chatRoom.hotline <- <-client.incoming
		}
	}()

	log.Printf("client %d joined\n", client.id)
}

// Listen waits for new messages and broadcasts them to all clients
func (chatRoom *ChatRoom) Listen() {
	go func() {
		for {
			chatRoom.Broadcast(<-chatRoom.hotline)
		}
	}()

	log.Println("chatroom listening...")
}

// Broadcast ...
func (chatRoom *ChatRoom) Broadcast(clientMsg ClientMessage) {
	for _, client := range chatRoom.clients {
		// Broadcast to all client except for the one that created the message
		if client.id != clientMsg.clientID {
			client.outgoing <- clientMsg.msg
		}
	}
}

// Client ...
type Client struct {
	conn     net.Conn
	id       int
	reader   *bufio.Reader
	writer   *bufio.Writer
	incoming chan ClientMessage
	outgoing chan string
}

// NewClient ...
func NewClient(conn net.Conn, id int) *Client {
	client := &Client{}
	client.conn = conn
	client.id = id
	client.reader = bufio.NewReader(conn)
	client.writer = bufio.NewWriter(conn)
	client.incoming = make(chan ClientMessage)
	client.outgoing = make(chan string)

	go client.Read()
	go client.Write()

	return client
}

// Read messages in from the client
func (client *Client) Read() {
	for {
		line, err := client.reader.ReadString('\n')
		if err != nil {
			// When a client disconnects the reader will return an EOF error.
			// When this happens, break and close the connection
			break
		}

		idStr := strconv.Itoa(client.id)
		msg := ClientMessage{
			clientID: client.id,
			msg:      idStr + ": " + line,
		}

		client.incoming <- msg
	}

	log.Printf("client %d disconneced\n", client.id)

	client.conn.Close()
}

// Write messages out to a client
func (client *Client) Write() {
	for msg := range client.outgoing {
		client.writer.WriteString(msg)
		client.writer.Flush()
	}
}

// ClientMessage ...
type ClientMessage struct {
	clientID int
	msg      string
}
