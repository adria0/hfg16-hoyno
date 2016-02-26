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
)

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	go h.run()

    router := gin.Default()
    router.LoadHTMLGlob("templates/*")

    router.GET("/", func(c *gin.Context) {
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
		c.HTML(200, "home.template",gin.H{
		    "whoami":name +" "+nif,
        })
	})
    router.GET("/ws", func(c *gin.Context) {
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
		Addr:      ":8443",
	}
	assert(server.ListenAndServeTLS("certs/localhost.pem", "certs/localhost.key"))
}



