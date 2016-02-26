package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "log"
)

func assert(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})

	})
    server := &http.Server{
        Handler:        router,
        Addr:           ":8443",
    }
    assert(server.ListenAndServeTLS("localhost.pem","localhost.key"))
}
