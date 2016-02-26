package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/boltdb/bolt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func assert(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var db *bolt.DB

func main() {

	var err error
	db, err = bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	go h.run()

	router := gin.Default()

	initWeb(router)

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
		Addr:      ":" + os.Getenv("PORT"),
	}
	assert(server.ListenAndServeTLS("certs/localhost.pem", "certs/localhost.key"))
}
