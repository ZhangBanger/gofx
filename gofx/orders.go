package gofx

import (
	"github.com/gin-gonic/gin"
	"time"
)

type Order struct {
	Timestamp int64   `db:"ts" json:"ts"`
	User      string  `json:"user" binding:"required" db:"user"`
	Security  string  `json:"security" binding:"required" db:"security"`
	Buy       bool    `json:"buy"`
	Quantity  uint32  `json:"quantity" binding:"required"`
	Price     float64 `json:"price" binding:"required"`
}

var orderChan chan Order

func MakeOrderChan(queueSize int) {
	orderChan = make(chan Order, queueSize)
}

func CreateOrder(c *gin.Context) {
	var order Order

	if c.Bind(&order) {
		order.Timestamp = time.Now().Unix()
		orderChan <- order
		c.JSON(202, gin.H{"status": "accepted", "payload": order})
	} else {
		c.JSON(400, gin.H{"status": "invalid fields"})
	}
}
