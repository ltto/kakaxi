package kakaxi

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"time"
)

var CA *x509.Certificate
var PK *rsa.PrivateKey

func init() {
	PKFilename := "kakaxi-ca-pk.pem"
	CAFilename := "kakaxi-ca-cert.pem"

	pemBytes, err := ioutil.ReadFile(PKFilename)
	if err != nil {
		if os.IsNotExist(err) {
			if PK, err = GeneratePK(2048); err != nil {
				panic(err)
			}
			if err := PKToFile(PK, PKFilename); err != nil {
				panic(err)
			}
			if pemBytes, err = ioutil.ReadFile(PKFilename); err != nil {
				panic(err)
			}
		}
	}
	PK, _ = LoadPKFromFile(pemBytes)
	if PK == nil {
		panic("PK ERR")
	}

	privateKeyData, err := ioutil.ReadFile(CAFilename)
	if err != nil {
		if os.IsNotExist(err) {
			var certByte []byte
			if CA, certByte, err = GenerateCA(PK, "MuYe", "KaKaXi", time.Now().AddDate(1, 0, 0)); err != nil {
				panic(err)
			}
			if err := CAToFile(CAFilename, certByte); err != nil {
				panic(err)
			}
			if privateKeyData, err = ioutil.ReadFile(CAFilename); err != nil {
				panic(err)
			}
		}
	}
	if CA, _ = LoadCertificateFromPEMBytes(privateKeyData); CA == nil {
		panic("CA ERR")
	}
}

const (
	PemHeaderPrivateKey  = "RSA PRIVATE KEY"
	PemHeaderCertificate = "CERTIFICATE"
)

func GeneratePK(bits int) (rsaKey *rsa.PrivateKey, err error) {
	rsaKey, err = rsa.GenerateKey(rand.Reader, bits)
	return
}
func PKToFile(rsaKey *rsa.PrivateKey, filename string) (err error) {
	keyOut, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %s", filename, err)
	}
	defer keyOut.Close()
	if err := pem.Encode(keyOut, &pem.Block{Type: PemHeaderPrivateKey, Bytes: x509.MarshalPKCS1PrivateKey(rsaKey)}); err != nil {
		return fmt.Errorf("unable to PEM encode private key: %s", err)
	}
	return
}

func GenerateCA(rsaKey *rsa.PrivateKey, organization string, name string, validUntil time.Time) (cert *x509.Certificate, certByte []byte, err error) {
	template := &x509.Certificate{
		SerialNumber: new(big.Int).SetInt64(int64(time.Now().UnixNano())),
		Subject: pkix.Name{
			Country:      []string{"CN"},
			Organization: []string{organization},
			CommonName:   name,
		},
		NotBefore:             time.Now().AddDate(0, -1, 0),
		NotAfter:              validUntil,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		Version:               3,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	certByte, err = x509.CreateCertificate(
		rand.Reader,       // secure entropy
		template,          // the template for the new cert
		template,          // cert that's signing this cert
		&rsaKey.PublicKey, // public key
		rsaKey,            // private key
	)
	if err != nil {
		return
	}
	cert, err = x509.ParseCertificate(certByte)
	if err != nil {
		return
	}
	return
}

func CAToFile(filename string, certByte []byte) (err error) {
	certOut, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to open %s for writing: %s", filename, err)
	}
	defer certOut.Close()
	return pem.Encode(certOut, &pem.Block{Type: PemHeaderCertificate, Bytes: certByte})
}
