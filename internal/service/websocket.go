package service

import (
	"bytes"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"log"
	"test.local/pkg/utils"
	"time"
)

const (
	writeWait = 10 *time.Second
	pongWait  = 60 * time.Second
	pingPeriod = (pongWait * 9) /10
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
	space = []byte{' '}
)

type Client struct {
	hub  *Hub
	conn *websocket.Conn
	send chan []byte
}

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

var vHub *Hub

func InitHub() {
	vHub = NewHub()
	go vHub.Run()
}
func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case message := <-h.broadcast:
			for client := range h.clients {
				select {
				case client.send <- message:
					log.Println("send mess", string(message))
				default:
					close(client.send)
					delete(h.clients, client)
					log.Println("run mess", message)
				}
			}

		}
	}
}

func NewClient() *Client {
	return &Client{
		hub:  vHub,
		send: make(chan []byte, 256),
	}
}

func (c *Client) Start(conn *websocket.Conn) {
	c.conn = conn
	c.hub.register <- c
	go c.writePump()
	go c.readPump()
}

func (c *Client) writePump(){
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for  {
		select {
		case messge, ok := <-c.send:
			log.Println("websocket write msg", string(messge))
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok{
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				utils.Log().Debug("service.Client.writePump return")
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil{
				utils.Log().Error("service.Client.writePump err", zap.Error(err))
				return
			}
			w.Write(messge)
			// 將剩餘的消息都推送出去
			n := len(c.send)
			for i:=0;i<n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil{
				utils.Log().Error("service.Client.writePump err", zap.Error(err))
				return
			}
		case <-ticker.C:
			utils.Log().Debug("service.Client.writePump ticker")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil{
				utils.Log().Error("service.Client.writePump err", zap.Error(err))
				return
			}
		}
	}
}

func(c *Client) readPump(){
	defer func() {
		c.hub.unregister <-c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		c.conn.WriteMessage(websocket.TextMessage, []byte(appData))
		return nil
	})

	for  {
		_, message, err := c.conn.ReadMessage()
		if err != nil{
			//if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure){
			//	utils.Log().Error("service.Client.readPump err", zap.Error(err))
			//}
			utils.Log().Error("service.Client.readPump err", zap.Error(err))
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		log.Println("websocket read msg", string(message))
		c.hub.broadcast <- message
	}
}