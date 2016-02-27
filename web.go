package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	//	"strings"
)

func buttonPressed(c *gin.Context, name string) bool {
	return c.DefaultPostForm(name, "undefined") != "undefined"
}

func initWeb(router *gin.Engine) {

	router.Static("/static", "web/static")
	router.LoadHTMLGlob("web/templates/*")

	router.GET("/chat", func(c *gin.Context) {
		cookie := getSessionUserId(c.Request)
		c.HTML(200, "chat.html", gin.H{
			"whoami":              cookie.nick,
			"IsAuthenticatedUser": !cookie.anonymous,
		})
	})

	router.GET("/", func(c *gin.Context) {
		if setSessionFromTLS(c) != Anonymous {
			c.Redirect(307, "/chat")
			return
		}
		setSessionFromAnonymous(c)
		c.HTML(200, "start.html", gin.H{})
	})

	router.GET("/info", func(c *gin.Context) {
		res := make(chan string)
		h.info <- res
		info := <-res
		c.String(200, "%v", info)
	})

	router.POST("/chatstatus/:active", func(c *gin.Context) {
		isActive := c.Param("active") == "on"
		ID := getSessionUserId(c.Request).ID
		setUserChatActivated(ID, isActive)
		c.String(200, "")
	})

	router.GET("/config", func(c *gin.Context) {
		ID := getSessionUserId(c.Request).ID
		user, err := load(ID)
		if err == ErrNotExists {
			user = User{
				ID:         ID,
				UserName:   ID,
				PublicName: "Falta nombre",
				ChatStatus: false,
			}
			save(ID, user)
			setUserChatActivated(ID, false)
		}

		fmt.Printf("%v", user)

		c.HTML(200, "config.html", gin.H{
			"user": user,
		})
	})

	router.POST("/config", func(c *gin.Context) {
		ID := getSessionUserId(c.Request).ID
		user := User{
			ID:          ID,
			UserName:    c.PostForm("UserName"),
			PublicName:  c.PostForm("PublicName"),
			Email:       c.PostForm("Email"),
			GroupName:   c.PostForm("GroupName"),
			GroupEmails: c.PostForm("GroupEmails"),
		}
		save(ID, user)
		user.ChatStatus = isUserChatActivated(ID)
		c.HTML(200, "config.html", gin.H{
			"user": user,
		})
	})
}
