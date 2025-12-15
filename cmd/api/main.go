package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": "1.0.0",
			"message": "TradeLog API is running",
		})
	})

	fmt.Println("🚀 Server running on port 8080")
	if err := r.Run(":8080"); err != nil {
		panic(err)
	}
}