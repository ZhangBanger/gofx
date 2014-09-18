package main

import (
	"flag"
	. "github.com/ZhangBanger/gofx/gofx"
	"github.com/gin-gonic/gin"
)

var queueSize int

func init() {
	flag.IntVar(&queueSize, "queue", 16, "size of queue")
	flag.Parse()
	InitDb()
	MakeOrderChan(queueSize)
}

func main() {
	go Process()

	r := gin.Default()
	r.POST("/orders", CreateOrder)
	r.GET("/orders", GetBook)

	r.Run(":8080")
}
