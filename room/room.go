package room

import (
	"bufio"
	"log"
	"net"

	"github.com/michaeljdennis/go-chat/message"
	"github.com/michaeljdennis/go-chat/room/user"
)

// Room ...
type Room struct {
	Users   map[string]*user.User
	Hotline chan Message
}

// New returns a new Room
func New() *Room {
	return &Room{
		Users:   make(map[string]*user.User),
		Hotline: make(chan Message),
	}
}

// Open ...
func (r *Room) Open() {
	r.CallHotline()
	r.Listen()
}

// CallHotline ...
func (r *Room) CallHotline() {
	log.Println("we're live on the hotline...")
	go func() {
		for {
			msg := <-r.Hotline
			for _, user := range r.Users {
				if user.ID != msg.UserID {
					log.Println("sending msg to", user.Name)
					writer := bufio.NewWriter(user.Conn)
					_, err := message.WriteMessage(writer, msg.Data)
					if err != nil {
						log.Println("Error broadcasting message:", err)
					}
				}
			}
		}
	}()
}

// Listen for new connections
func (r *Room) Listen() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err)
	}
	defer ln.Close()

	for {
		// Wait for connection
		conn, err := ln.Accept()
		if err != nil {
			log.Fatalln("Error accepting: ", err.Error())
		}

		// Create new user
		u := user.New(conn)

		// Add user to room
		r.Users[u.ID] = u

		// Handle user connection
		go r.HandleConnection(u)
	}
}

// HandleConnection ...
func (r *Room) HandleConnection(user *user.User) {
	defer user.Conn.Close()

	// Remove user after disconnect
	defer delete(r.Users, user.ID)

	var err error

	// Set username
	reader := bufio.NewReader(user.Conn)
	user.Name, err = message.ReadMessage(reader)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("user %s (%s) connected\n", user.Name, user.ID)

	// Send UUID acknowledgment
	writer := bufio.NewWriter(user.Conn)
	_, err = message.WriteMessage(writer, user.ID)
	if err != nil {
		log.Fatalln(err)
	}

	// Wait for new messages
	for {
		// Wait for new string to be read
		reader := bufio.NewReader(user.Conn)
		msg, err := message.ReadMessage(reader)
		if err != nil {
			// TODO: differentiate between EOF and all other errors
			log.Printf("%s disconnected\n", user.Name)
			return
		}

		// Log user messages #NSA
		log.Printf("%s: %s", user.Name, msg)

		// Send message to hotline
		log.Println("sending message to hotline:", msg)
		r.Hotline <- Message{
			Data:   user.Name + ": " + msg,
			UserID: user.ID,
		}
	}
}

// Message ...
type Message struct {
	Data   string
	UserID string
}
