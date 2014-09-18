package main

import "github.com/gin-gonic/gin"

type Order struct {
	User     string  `json:"user" binding:"required"`
	Security string  `json:"security" binding:"required"`
	Buy      bool    `json:"buy"`
	Quantity uint32  `json:"quantity" binding:"required"`
	Price    float64 `json:"price" binding:"required"`
}

func main() {
	r := gin.Default()
	r.POST(
		"/orders",
		func(c *gin.Context) {
			var order Order

			if c.Bind(&order) {
				c.JSON(202, gin.H{"status": "accepted", "payload": order})
			} else {
				c.JSON(400, gin.H{"status": "invalid fields"})
			}

		},
	)

	r.Run(":8080")
}
