package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strings"
)

func buttonPressed(c *gin.Context, name string) bool {
	return c.DefaultPostForm(name, "undefined") != "undefined"
}

func initWeb(router *gin.Engine) {

	router.Static("/static", "web/static")
	router.LoadHTMLGlob("web/templates/*")

	router.GET("/chat", func(c *gin.Context) {
		var (
			name string
			nif  string
		)
		if len(c.Request.TLS.PeerCertificates) > 0 {
			cn := c.Request.TLS.PeerCertificates[0].Subject.CommonName
			split := strings.Split(cn, " - NIF ")
			name = split[0]
			nif = split[1]
		}
		c.HTML(200, "chat.html", gin.H{
			"whoami": name + " " + nif,
		})
	})

	router.GET("/", func(c *gin.Context) {
		setSessionUserId(c)
		c.HTML(200, "start.html", gin.H{})
	})

	router.GET("/info", func(c *gin.Context) {
		res := make(chan string)
		h.info <- res
		info := <-res
		c.String(200, "%v", info)
	})

	router.POST("/config/chatstatus/:active", func(c *gin.Context) {
		isActive := c.Param("active") == "on"
		ID := getSessionUserId(c)
		setChatStatus(ID, isActive)
		c.String(200, "")
	})

	router.GET("/config", func(c *gin.Context) {
		ID := getSessionUserId(c)
        user, err := load(ID)
        if err == ErrNotExists {
				user = User{
					ID:         ID,
					UserName:   ID,
					PublicName: "Falta nombre",
				}
				save(ID, user)
                setChatStatus(ID,false)
        }

		fmt.Printf("%v", user)

		c.HTML(200, "config.html", gin.H{
			"user": user,
		})
	})

	router.POST("/config", func(c *gin.Context) {
		ID := getSessionUserId(c)
		user := User{
			ID:          ID,
			UserName:    c.PostForm("UserName"),
			PublicName:  c.PostForm("PublicName"),
			Email:       c.PostForm("Email"),
			GroupName:   c.PostForm("GroupName"),
			GroupEmails: c.PostForm("GroupEmails"),
		}
		save(ID, user)

		c.HTML(200, "config.html", gin.H{
			"user": user,
		})
	})
}
