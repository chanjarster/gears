/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package confs

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"strings"
	"time"
)

// Config keys:
//  | Environment            |  Flag                   |  Description                              |
//  |------------------------|-------------------------|-------------------------------------------|
//  | CERT_CHAIN             | -cert-chain             | Path to cert when communicate client auth |
//  |                        |                         | enabled origin server, including          |
//  |                        |                         | intermediate certs. x.509 pem encoded.    |
//  | PRIVATE_KEY            | -private-key            | Path to private key when communicate      |
//  |                        |                         | client auth enabled origin server,        |
//  |                        |                         | including intermediate certs.             |
//  |                        |                         | PKCS #1 pem encoded.                      |
//  | SSL_TRUST_MODE         | -ssl-trust-mode         | Trust mode when verifying server          |
//  |                        |                         | certificates. Available modes are:        |
//  |                        |                         | OS: use host root CA set.                 |
//  |                        |                         | INSECURE: do not verify.                  |
//  |                        |                         | CUSTOM: verify by custom cert.            |
//  | SSL_TRUST_CERTS        | -ssl-trust-certs        | Path to certificates that clients use     |
//  |                        |                         | when verifying server certificates, only  |
//  |                        |                         | useful when client-ssl-trust-mode set     |
//  |                        |                         | to CUSTOM. X.509 PEM encoded              |
//  | READ_TIMEOUT           | -read-timeout           |                                           |
//  | WRITE_TIMEOUT          | -write-timeout          |                                           |
//  | MAX_CONN_DURATION      | -max-conn-duration      |                                           |
//  | MAX_CONNS_PER_HOST     | -max-conns-per-host     |                                           |
//  | MAX_IDLE_CONN_DURATION | -max-idle-conn-duration |                                           |
//  | MAX_CONN_WAIT_TIMEOUT  | -max-conn-wait-timeout  |                                           |
// Note: if FastHttpClientConf is nested in another struct, add corresponding prefix.
type FastHttpClientConf struct {
	CertChain     string // Path to cert when communicate client auth enabled origin server, including intermediate certs. x.509 pem encoded.
	PrivateKey    string // Path to private key when communicate client auth enabled origin server, including intermediate certs. PKCS #1 pem encoded.
	SslTrustMode  string // Trust mode when verifying server certificates. Available modes are OS: use host root CA set. INSECURE: do not verify. CUSTOM: verify by custom cert.
	SslTrustCerts string // Path to certificates that clients use when verifying server certificates, only useful when client-ssl-trust-mode set to CUSTOM. X.509 PEM encoded

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
		errLogger.Fatalf("load cert %s failed\n", file)
	}
}
