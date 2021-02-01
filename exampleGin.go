package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
	//默认端口：8080
	r := gin.Default()
	r.GET("/exampleGin", exampleGin)

	r.Run()
}

func exampleGin(c *gin.Context) {
	id := c.Query("id")

	c.String(http.StatusOK, "Hello %s", id)
}
