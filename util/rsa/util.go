package rsautil

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
)

func MustReadPrivateKey(pem string) *rsa.PrivateKey {
	p, err := ReadPrivateKey(pem)
	if err != nil {
		panic(err)
	}
	return p
}

func ReadPrivateKey(pem string) (*rsa.PrivateKey, error) {
	pem = normalizePem(pem, "RSA PRIVATE KEY")
	pemByte := []byte(pem)
	return ReadPrivateKeyBytes(pemByte)
}

func MustReadPrivateKeyBytes(pemBytes []byte) *rsa.PrivateKey {
	p, err := ReadPrivateKeyBytes(pemBytes)
	if err != nil {
		panic(err)
	}
	return p
}

func ReadPrivateKeyBytes(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid RSA PRIVATE KEY")
	}
	if block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("expected type: RSA PRIVATE KEY, got: " + block.Type)
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func MustReadPublicKey(pem string) *rsa.PublicKey {
	p, err := ReadPublicKey(pem)
	if err != nil {
		panic(err)
	}
	return p
}

func ReadPublicKey(pem string) (*rsa.PublicKey, error) {
	pem = normalizePem(pem, "PUBLIC KEY")
	pemByte := []byte(pem)
	return ReadPublicKeyBytes(pemByte)
}

func MustReadPublicKeyBytes(pemBytes []byte) *rsa.PublicKey {
	p, err := ReadPublicKeyBytes(pemBytes)
	if err != nil {
		panic(err)
	}
	return p
}

func ReadPublicKeyBytes(pemBytes []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, errors.New("invalid PUBLIC KEY")
	}
	if block.Type != "PUBLIC KEY" {
		return nil, errors.New("expected type: RSA PUBLIC KEY, got: " + block.Type)
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub.(*rsa.PublicKey), nil
}

func normalizePem(pemStr string, typ string) string {
	if strings.Index(pemStr, "-----BEGIN") == -1 || strings.Index(pemStr, "-----END") == -1 {
		pemStr = fmt.Sprintf("-----BEGIN %s-----\n%s\n-----END %s-----", typ, pemStr, typ)
	}
	return pemStr
}

func ExtractPublicKeyPem(privateKey *rsa.PrivateKey) (string, error) {
	bs, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		return "", err
	}
	b := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: bs,
	}
	buf := &bytes.Buffer{}
	err = pem.Encode(buf, b)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}
