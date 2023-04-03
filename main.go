package main

// * Template file for testing API. NOT EDIT

import (
	wwz "api/wwz/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	api := r.Group("/api/wwz")
	wwz.RegisterRoutes(api)

	r.Run(":3000")
}
