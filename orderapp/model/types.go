package model

import (
	"github.com/ServiceWeaver/weaver"
)

const (
	Email string = "email"
	SMS   string = "sms"
)

type OrderPayment struct {
	weaver.AutoMarshal
	OrderID    string
	CustomerID string
	Amount     float64
}

type Notification struct {
	weaver.AutoMarshal
	OrderID    string
	CustomerID string
	Event      string
	Modes      []string
}
