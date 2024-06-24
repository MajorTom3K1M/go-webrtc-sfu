package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

type Client struct {
	Hub     *Hub
	Conn    *websocket.Conn
	Send    chan Message
	ID      string
	Channel string
	sync.Mutex
}

type Hub struct {
	Clients       map[*Client]bool
	Broadcast     chan Message
	Register      chan *Client
	Unregister    chan *Client
	Channels      map[string]map[*Client]bool
	PeerChannels  map[string][]PeerConnectionState
	TrackChannels map[string]map[string]*webrtc.TrackLocalStaticRTP
	sync.RWMutex
}

func (c *Client) WriteJSON(v interface{}) error {
	c.Lock()
	defer c.Unlock()
	return c.Conn.WriteJSON(v)
}

func (c *Client) readPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	ps, err := NewPeerConnectionState(c)
	if err != nil {
		log.Println(err)
		return
	}

	defer ps.peerConnection.Close()

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("WebSocket read message error:", err)
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("JSON unmarshal error:", err)
			continue
		}

		fmt.Printf("Received : %s \n", msg.Type)
		switch msg.Type {
		case "answer":
			if err := ps.peerConnection.SetRemoteDescription(msg.Answer); err != nil {
				log.Println("Failed to set remote description:", err)
			}
		case "candidate":
			if err := ps.peerConnection.AddICECandidate(msg.Candidate); err != nil {
				log.Println("Failed to add ICE candidate:", err)
			}
		}
	}
}

func (c *Client) writePump() {
	for message := range c.Send {
		c.Lock()
		err := c.Conn.WriteJSON(message)
		c.Unlock()
		if err != nil {
			log.Println("WebSocket write error:", err)
			return
		}
	}
	c.Conn.Close()
}

func newWebSocketHub() *Hub {
	return &Hub{
		Clients:       make(map[*Client]bool),
		Broadcast:     make(chan Message),
		Register:      make(chan *Client),
		Unregister:    make(chan *Client),
		Channels:      make(map[string]map[*Client]bool),
		PeerChannels:  make(map[string][]PeerConnectionState),
		TrackChannels: make(map[string]map[string]*webrtc.TrackLocalStaticRTP),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.Channels[client.Channel]; !ok {
				h.Channels[client.Channel] = make(map[*Client]bool)
			}
			h.Channels[client.Channel][client] = true
			h.Clients[client] = true
		case client := <-h.Unregister:
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)
				if clients, ok := h.Channels[client.Channel]; ok {
					delete(clients, client)
					if len(clients) == 0 {
						delete(h.Channels, client.Channel)
					}
				}
			}
		case message := <-h.Broadcast:
			if clients, ok := h.Channels[message.Channel]; ok {
				for client := range clients {
					select {
					case client.Send <- message:
					default:
						close(client.Send)
						delete(h.Clients, client)
					}
				}
			}
		}
	}
}
