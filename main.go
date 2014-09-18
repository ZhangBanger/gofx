package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.POST(
		"/orders",
		func(c *gin.Context) {
			c.String(202, "created")
		},
	)

	r.Run(":8080")
}
