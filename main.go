package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
    "log"
    "crypto/x509"
    "crypto/tls"
    "io/ioutil"
    "fmt"
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
    caCertPool := x509.NewCertPool()
    fileInfos, err := ioutil.ReadDir("certs/ca")
    assert(err)
    for _, fileInfo := range fileInfos {
        if !fileInfo.IsDir() {
            caCert, err := ioutil.ReadFile("certs/ca/"+fileInfo.Name())
            assert(err)
            if (!caCertPool.AppendCertsFromPEM(caCert)){
                assert(fmt.Errorf("Failed to read ca %v", fileInfo.Name))
            }
        }
    }

    tlsConfig := &tls.Config{
        // Reject any TLS certificate that cannot be validated
        ClientAuth: tls.RequireAndVerifyClientCert,
        // Ensure that we only use our "CA" to validate certificates
        ClientCAs: caCertPool,
        // PFS because we can but this will reject client with RSA certificates
//        CipherSuites: []uint16{tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384},
        // Force it server side
        PreferServerCipherSuites: true,
        // TLS 1.2 because we can
        MinVersion: tls.VersionTLS12,
    }


    server := &http.Server{
        Handler:        router,
        TLSConfig:       tlsConfig,
        Addr:           ":8443",
    }
    assert(server.ListenAndServeTLS("certs/localhost.pem","certs/localhost.key"))
}
