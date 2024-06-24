package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pion/webrtc/v3"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	Type      string                    `json:"type"`
	Name      string                    `json:"name,omitempty"`
	Offer     webrtc.SessionDescription `json:"offer,omitempty"`
	Answer    webrtc.SessionDescription `json:"answer,omitempty"`
	Candidate webrtc.ICECandidateInit   `json:"candidate,omitempty"`
	Channel   string                    `json:"channel,omitempty"`
}

func handleWebsocket(w http.ResponseWriter, r *http.Request, hub *Hub) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	channel := r.URL.Query().Get("channel")
	if channel == "" {
		log.Print("No channel specified")
		return
	}

	client := &Client{
		Hub:     hub,
		Conn:    conn,
		Send:    make(chan Message, 256),
		Channel: channel,
	}
	client.Hub.Register <- client

	go client.readPump()
	go client.writePump()
}
