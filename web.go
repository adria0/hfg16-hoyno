package main

import (
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"strings"
)

var (
	UserNameKey    = []byte{1}
	PublicNameKey  = []byte{2}
	EmailKey       = []byte{3}
	GroupNameKey   = []byte{4}
	GroupEmailsKey = []byte{5}
)

type User struct {
	ID          string
	UserName    string
	PublicName  string
	Email       string
	GroupName   string
	GroupEmails string
}

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

	router.GET("/config", func(c *gin.Context) {
		ID := getSessionUserId(c)
		var user User

		db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(ID))
			if b != nil {
				user = User{
					ID:          ID,
					UserName:    string(b.Get(UserNameKey)),
					PublicName:  string(b.Get(PublicNameKey)),
					Email:       string(b.Get(EmailKey)),
					GroupName:   string(b.Get(GroupNameKey)),
					GroupEmails: string(b.Get(GroupEmailsKey)),
				}
			} else {
				user = User{
					ID:         ID,
					UserName:   ID,
					PublicName: "Falta nombre",
				}
			}
			return nil
		})

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

		db.Update(func(tx *bolt.Tx) error {
			b, err := tx.CreateBucketIfNotExists([]byte(ID))
			if err != nil {
				return err
			}
			b.Put(UserNameKey, []byte(user.UserName))
			b.Put(PublicNameKey, []byte(user.PublicName))
			b.Put(EmailKey, []byte(user.Email))
			b.Put(GroupNameKey, []byte(user.GroupName))
			b.Put(GroupEmailsKey, []byte(user.GroupEmails))
			return nil
		})

		c.HTML(200, "config.html", gin.H{
			"user": user,
		})
	})
}
