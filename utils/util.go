package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"crypto/tls"
	"crypto/x509"
)

var (
	servertls12client string
	serverfileca      string
)

const (
	SERVERPREFIX = "server"
)

type ServerConfig struct {

	// ReadTimeout is the maximum duration for reading the entire
	// request, including the body.
	//
	// Because ReadTimeout does not let Handlers make per-request
	// decisions on each request body's acceptable deadline or
	// upload rate, most users will prefer to use
	// ReadHeaderTimeout. It is valid to use them both.
	ReadTimeout time.Duration

	// ReadHeaderTimeout is the amount of time allowed to read
	// request headers. The connection's read deadline is reset
	// after reading the headers and the Handler can decide what
	// is considered too slow for the body.
	ReadHeaderTimeout time.Duration `env:"rhtimeout"`

	// WriteTimeout is the maximum duration before timing out
	// writes of the response. It is reset whenever a new
	// request's header is read. Like ReadTimeout, it does not
	// let Handlers make decisions on a per-request basis.
	WriteTimeout time.Duration

	// IdleTimeout is the maximum amount of time to wait for the
	// next request when keep-alives are enabled. If IdleTimeout
	// is zero, the value of ReadTimeout is used. If both are
	// zero, ReadHeaderTimeout is used.
	IdleTimeout time.Duration

	// MaxHeaderBytes controls the maximum number of bytes the
	// server will read parsing the request header's keys and
	// values, including the request line. It does not limit the
	// size of the request body.
	// If zero, DefaultMaxHeaderBytes is used.
	MaxHeaderBytes int `env:"maxbytes"`

	Serverfileca         string `env:"FILECA"`
	Serverfileprivatekey string `env:"PRIVATEKEY"`
	Serverfilepubkey     string `env:"PUBLICKEY"`
	Servertls12client    string `env:"TLS12STATUS" envDefault:"OFF"`
}

func GetServerConfig() *ServerConfig {
	var serverCfg ServerConfig
	return &serverCfg
}

func GetServerTlsConfig() *tls.Config {
	if servertls12client == "ON" {

		caCert, err := ioutil.ReadFile(serverfileca)
		if err != nil {
			fmt.Println("Error : ", err)

		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		tlsConfig := &tls.Config{
			ClientCAs:  caCertPool,
			ClientAuth: tls.RequireAndVerifyClientCert,
		}
		tlsConfig.BuildNameToCertificate()
		return tlsConfig
	}
	return &tls.Config{InsecureSkipVerify: true}

}

// Simple helper function to read an environment or return a default value
func GetEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}
