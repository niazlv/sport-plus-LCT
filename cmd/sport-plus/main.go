package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/niazlv/sport-plus-LCT/internal/routes"
)

func main() {
	app := gin.Default()

	app.Use(cors.Default())

	routes.Setup(app)

	app.Run(":8000")
}
