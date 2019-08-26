package message

import (
	"bufio"
	"log"
	"strings"
)

// ReadMessage ...
func ReadMessage(r *bufio.Reader) (string, error) {
	msg, err := r.ReadString('\n')
	if err != nil {
		// log.Println("ReadMessage", err)
		return "", err
	}
	msg = strings.TrimSuffix(msg, "\n")
	return msg, nil
}

// WriteMessage ...
func WriteMessage(w *bufio.Writer, msg string) (int, error) {
	n, err := w.WriteString(msg + "\n")
	if err != nil {
		log.Println("WriteMessage", err)
		return 0, err
	}
	w.Flush()
	return n, nil
}
