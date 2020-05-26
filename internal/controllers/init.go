package controllers

type Controllers struct {
	User *User
}

var Ctrl *Controllers

func Init() {
	Ctrl = &Controllers{
		User: NewUser(),
	}
}
