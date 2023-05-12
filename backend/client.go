package main

import (
	"strings"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub *Hub
	conn *websocket.Conn
	send chan []byte
	name string
}

type Message struct {
	Sender string `json:"sender"`
	Content string `json:"body"`
	Recipient string `json:"recipient"`
	Newname string `json:"newname"`
	Action string `json:"action"`
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Printf("error: %v", err)
			}
			break
		}

		msg := Message{}
		err = json.Unmarshal(message, &msg)
		if err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		msg.Sender = c.name

		if strings.HasPrefix(msg.Content, "/chname") {
			parts := strings.SplitN(msg.Content, " ", 2)
			if len(parts) == 2 {
				if c.name == parts[1] {
					msg.Content = "You cannot change name to the original."
				} else if parts[1] == "anonymous" {
					msg.Content = "You cannot change name to 'anonymous'."
				} else {
					msg.Content = c.name + " change name to " + parts[1]
					c.name = parts[1]
					msg.Newname = c.name
				}
			}
			msg.Action = "chname"
		} else if strings.HasPrefix(msg.Content, "/to") {
			parts := strings.SplitN(msg.Content, " ", 3)
			if len(parts) == 3 {
				msg.Recipient = parts[1]
				msg.Content = parts[2]
			}
			msg.Action = "to"
		}

		message, err = json.Marshal(msg)
		if err != nil {
			log.Printf("Error marshaling message: %v", err)
			continue
		}

		c.hub.broadcast <- message
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256), name: "anonymous"}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}