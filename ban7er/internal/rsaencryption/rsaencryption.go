package utils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func GenerateKeyPair() (privKeyPEM string, pubKeyPEM string, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("error generating RSA key pair: %v", err)
	}

	privKeyPEM, pubKeyPEM = keysToPEM(privateKey, &privateKey.PublicKey)

	return privKeyPEM, pubKeyPEM, nil
}

func keysToPEM(
	privKey *rsa.PrivateKey,
	pubKey *rsa.PublicKey,
) (privKeyPEM string, pubKeyPEM string) {
	privKeyBytes := x509.MarshalPKCS1PrivateKey(privKey)
	privKeyPEM = string(
		pem.EncodeToMemory(
			&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: privKeyBytes,
			},
		),
	)

	pubKeyBytes := x509.MarshalPKCS1PublicKey(pubKey)
	pubKeyPEM = string(
		pem.EncodeToMemory(
			&pem.Block{
				Type:  "PUBLIC KEY",
				Bytes: pubKeyBytes,
			},
		),
	)

	return privKeyPEM, pubKeyPEM
}
