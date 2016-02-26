package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
    "os"
    "github.com/boltdb/bolt"
)

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var db *bolt.DB

func main() {

    db, err := bolt.Open("my.db", 0600, nil)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    go h.run()

    router := gin.Default()
    router.Static("/static", "web/static")
    router.LoadHTMLGlob("web/templates/*")


    router.GET("/chat", func(c *gin.Context) {
        var (
            name string
            nif string
        )
        if len(c.Request.TLS.PeerCertificates)>0 {
            cn := c.Request.TLS.PeerCertificates[0].Subject.CommonName
		    split := strings.Split(cn, " - NIF ")
		    name = split[0]
		    nif = split[1]
        }
		c.HTML(200, "chat.html",gin.H{
		    "whoami":name +" "+nif,
        })
	})

    router.GET("/", func(c *gin.Context) {
		c.HTML(200, "start.html",gin.H{
        })
	})

    router.GET("/info", func(c *gin.Context) {
        res := make (chan string)
        h.info <- res
        info := <- res
        c.String(200, "%v",info)
	})

    router.GET("/chatws", func(c *gin.Context) {
        serveWs(c.Writer, c.Request)
    })

    caCertPool := x509.NewCertPool()
	fileInfos, err := ioutil.ReadDir("certs/ca")
	assert(err)
	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			caCert, err := ioutil.ReadFile("certs/ca/" + fileInfo.Name())
			assert(err)
			if !caCertPool.AppendCertsFromPEM(caCert) {
				assert(fmt.Errorf("Failed to read ca %v", fileInfo.Name))
			}
		}
	}

	tlsConfig := &tls.Config{
		ClientAuth:               tls.RequestClientCert,
		ClientCAs:                caCertPool,
		PreferServerCipherSuites: true,
		MinVersion:               tls.VersionTLS12,
	}

	server := &http.Server{
		Handler:   router,
		TLSConfig: tlsConfig,
		Addr:      ":"+os.Getenv("PORT"),
	}
	assert(server.ListenAndServeTLS("certs/localhost.pem", "certs/localhost.key"))
}



