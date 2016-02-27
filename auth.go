package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
	//    "log"
)

const (
	Anonymous = ""
)

type Cookie struct {
	expires   int64
	ID        string
	anonymous bool
	nick      string
    team      string
}

var (
	cookies = map[string]*Cookie{}
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

func getSessionUserId(r *http.Request) *Cookie {

	validSessionUserId := func(token string) *Cookie {
		now := time.Now().Unix()

		mutex.Lock()
		defer mutex.Unlock()

		if cookie, exists := cookies[token]; exists {
			if now > cookie.expires {
				delete(cookies, token)
				return nil
			}
			return cookie
		}
		return nil

	}

	for _, cookie := range r.Cookies() {
		if cookie.Name == "token" {
			return validSessionUserId(cookie.Value)
		}

	}
	return nil
}

func setSessionFromAnonymous(c *gin.Context) string {
	token128 := fmt.Sprintf("%x%x%x%x",
		rand.Uint32(), rand.Uint32(), rand.Uint32(), rand.Uint32())
	expires := time.Now().Unix() + 7*24*3600

	ID := "Random_user"

	mutex.Lock()
	cookies[token128] = &Cookie{expires: expires, ID: ID, anonymous: true, nick: ID}
	mutex.Unlock()

	cookie := http.Cookie{Name: "token", Value: token128}
	http.SetCookie(c.Writer, &cookie)

	return ID
}

func setSessionFromTLS(c *gin.Context) string {

	if len(c.Request.TLS.PeerCertificates) == 0 {
		fmt.Println("------ user ANONIM ------")
		return Anonymous
	}

	cn := c.Request.TLS.PeerCertificates[0].Subject.CommonName
	split := strings.Split(cn, " - NIF ")
	ID := "NIF" + split[1]
	nick := ID

	token128 := fmt.Sprintf("%x%x%x%x",
		rand.Uint32(), rand.Uint32(), rand.Uint32(), rand.Uint32())
	expires := time.Now().Unix() + 7*24*3600

	user, err := load(ID)
	if err == nil {
		fmt.Printf("------ user: %v ------", user)
		nick = user.PublicName
	} else {

		fmt.Printf("------ user: %v ------", err)
	}

	mutex.Lock()
	cookies[token128] = &Cookie{expires: expires, ID: ID, anonymous: false, nick: nick}
	mutex.Unlock()

	cookie := http.Cookie{Name: "token", Value: token128}
	http.SetCookie(c.Writer, &cookie)

	return ID
}
