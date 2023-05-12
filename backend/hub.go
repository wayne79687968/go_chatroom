package main

import (
	"encoding/json"
	"log"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			var msg Message
			err := json.Unmarshal(message, &msg)
			if err != nil {
				log.Printf("Error unmarshaling message: %v", err)
				continue
			}
			if "" == msg.Sender {
				for client := range h.clients {
					if client.name == msg.Sender {
						select {
						case client.send <- message:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				}
			} else if "" == msg.Recipient {
				for client := range h.clients {
					select {
					case client.send <- message:
					default:
						close(client.send)
						delete(h.clients, client)
					}
				}
			} else {
				for client := range h.clients {
					if client.name == msg.Recipient || client.name == msg.Sender {
						select {
						case client.send <- message:
						default:
							close(client.send)
							delete(h.clients, client)
						}
					}
				}
			}
		}
	}
}