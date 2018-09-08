package chessmatch

import "github.com/gorilla/websocket"

// Client represents a client who asked for a connection
type Client struct {
	Hub *Hub

	//the match name
	Match string

	// websocket connection
	Conn *websocket.Conn

	// channel of outbund movements
	Send chan []byte
}

func (p *Client) Read() {
	// Starts webscoket connection

	// close of the connection after the end of function

	// loop on the data

	// pass the data trought the channel to the player
}

func (p *Client) Write() {
	// Starts a websocket connection

	// process the data to be sent over the socket

	// send the data to the player's room
}
