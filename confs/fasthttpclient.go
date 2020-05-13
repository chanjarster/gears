package confs

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type FastHttpClientConf struct {
	CertChain     string // cert when communicate client auth enabled origin server, including intermediate certs. x.509 pem encoded.
	PrivateKey    string // private key when communicate client auth enabled origin server, including intermediate certs. PKCS #1 pem encoded.
	SslTrustMode  string // trust mode when verifying server certificates. Available modes are OS: use host root CA set. INSECURE: do not verify. CUSTOM: verify by custom cert.
	SslTrustCerts string // certificate that clients use when verifying server certificates, only useful when client-ssl-trust-mode set to CUSTOM. X.509 PEM encoded

	ReadTimeout         time.Duration
	WriteTimeout        time.Duration
	MaxConnDuration     time.Duration
	MaxConnsPerHost     int
	MaxIdleConnDuration time.Duration
	MaxConnWaitTimeout  time.Duration
}

func (c *FastHttpClientConf) String() string {
	return fmt.Sprintf("{CertChain: %s, PrivateKey: %s, SslTrustMode: %s, SslTrustCerts: %s, ReadTimeout:%s, WriteTimeout: %s, MaxConnDuration: %s, MaxConnsPerHost: %d, MaxIdleConnDuration: %s, MaxConnWaitTimeout: %s}",
		c.CertChain, c.PrivateKey, c.SslTrustMode, c.SslTrustCerts,
		c.ReadTimeout, c.WriteTimeout, c.MaxConnDuration, c.MaxConnsPerHost,
		c.MaxIdleConnDuration, c.MaxConnWaitTimeout)
}

type FastHttpClientCustomizer func(hc *fasthttp.Client)

func NewFastHttpClient(hcConf *FastHttpClientConf, customizer FastHttpClientCustomizer) *fasthttp.Client {

	tlsConfigured := false
	tlsConfig := &tls.Config{}

	if hcConf.CertChain != "" || hcConf.PrivateKey != "" {
		keyPair, err := tls.LoadX509KeyPair(hcConf.CertChain, hcConf.PrivateKey)
		if err != nil {
			panic(err)
		}
		tlsConfig.Certificates = []tls.Certificate{keyPair}
		tlsConfigured = true
	}

	if hcConf.SslTrustMode == "CUSTOM" && hcConf.SslTrustCerts != "" {
		certPool := x509.NewCertPool()
		for _, cert := range strings.Split(hcConf.SslTrustCerts, ",") {
			appendCertPool(certPool, cert)
		}
		tlsConfig.RootCAs = certPool
	} else if hcConf.SslTrustMode == "INSECURE" {
		tlsConfig.InsecureSkipVerify = true
		tlsConfigured = true
	}

	client := &fasthttp.Client{
		MaxConnsPerHost:     hcConf.MaxConnsPerHost,
		MaxIdleConnDuration: hcConf.MaxIdleConnDuration,
		MaxConnDuration:     hcConf.MaxConnDuration,
		ReadTimeout:         hcConf.ReadTimeout,
		WriteTimeout:        hcConf.WriteTimeout,
		MaxConnWaitTimeout:  hcConf.MaxConnWaitTimeout,
	}

	if tlsConfigured {
		client.TLSConfig = tlsConfig
	}
	if customizer != nil {
		customizer(client)
	}

	return client
}

func appendCertPool(certPool *x509.CertPool, file string) {
	cert, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	if !certPool.AppendCertsFromPEM(cert) {
		log.Fatalf("load cert %s failed\n", file)
	}
}
