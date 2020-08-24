package websocket

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type person struct {
	Id       string `json:"_id"`
	Insultos int    `json:"insultos"`
}

func (c *Client) Read() {
	defer func() {
		c.Pool.Unregister <- c
		_ = c.Conn.Close()
	}()

	for {
		_, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		var update person
		_ = json.Unmarshal(p, &update)

		message := person{Id: update.Id, Insultos: update.Insultos}
		c.Pool.Broadcast <- message
		fmt.Printf("person Received: %+v\n", message)
	}
}
