package kakaxi

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
)

func LoadCertificateFromPEMBytes(pemBytes []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(pemBytes)
	if block == nil {
		return nil, fmt.Errorf("unable to decode PEM encoded certificate")
	}
	certificate, err := x509.ParseCertificate(block.Bytes)
	return certificate, err
}

func LoadPKFromFile(privateKeyData []byte) (key *rsa.PrivateKey, err error) {
	block, _ := pem.Decode(privateKeyData)
	if block == nil {
		return nil, fmt.Errorf("unable to decode PEM encoded private key data: %s", err)
	}
	return x509.ParsePKCS1PrivateKey(block.Bytes)
}
