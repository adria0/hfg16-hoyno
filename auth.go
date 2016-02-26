package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	Anonymous = ""
)

type Cookie struct {
	expires int64
	ID      string
}

var (
	cookies = map[string]Cookie{}
	mutex   = &sync.Mutex{}
)

func initAuth() {
	ticker := time.NewTicker(time.Hour)
	go func() {
		for _ = range ticker.C {
			now := time.Now().Unix()

			mutex.Lock()

			for k, cookie := range cookies {
				if now > cookie.expires {
					delete(cookies, k)
					fmt.Println("Removed", k)

				}

			}
			mutex.Unlock()

		}

	}()

}

func getSessionUserId(c *gin.Context) string {

	validSessionUserId := func(token string) string {
		now := time.Now().Unix()

		mutex.Lock()
		defer mutex.Unlock()

		if cookie, exists := cookies[token]; exists {
			if now > cookie.expires {
				delete(cookies, token)
				return Anonymous
			}
			return cookie.ID
		}
		return Anonymous

	}

	for _, cookie := range c.Request.Cookies() {
		if cookie.Name == "token" {
			return validSessionUserId(cookie.Value)
		}

	}
	return Anonymous
}

func setSessionUserId(c *gin.Context) string {

	if len(c.Request.TLS.PeerCertificates) == 0 {
		return Anonymous
	}

	cn := c.Request.TLS.PeerCertificates[0].Subject.CommonName
	split := strings.Split(cn, " - NIF ")
	ID := "NIF" + split[1]

	token128 := fmt.Sprintf("%x%x%x%x",
		rand.Uint32(), rand.Uint32(), rand.Uint32(), rand.Uint32())
	expires := time.Now().Unix() + 7*24*3600

	mutex.Lock()
	cookies[token128] = Cookie{expires: expires, ID: ID}
	mutex.Unlock()

	cookie := http.Cookie{Name: "token", Value: token128}
	http.SetCookie(c.Writer, &cookie)

	return ID

}
