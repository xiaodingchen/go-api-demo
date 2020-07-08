package controllers

type Controllers struct {
	User *User
	Ws *WebSocket
}

var Ctrl *Controllers

func Init() {
	Ctrl = &Controllers{
		User: NewUser(),
		Ws: NewWs(),
	}
}
