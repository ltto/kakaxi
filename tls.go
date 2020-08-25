package kakaxi

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"net"
)

func OnTLS(accept net.Conn, host string) (err error) {
	_, _ = accept.Write([]byte("HTTP/1.1 200 Connection Established\r\n\r\n"))
	config := tls.Config{GetCertificate: TLSHandshake}
	conn := tls.Server(accept, &config)
	if err = conn.Handshake(); err != nil {
		return err
	}
	header, body, paths, method := DumpRequest(conn)
	doHeader, bodyB, err := ProxyHTTP("https://"+host+paths, method, header, body)
	if err != nil {
		return err
	}
	Writer(conn, doHeader, bodyB)
	conn.Close()
	return
}

func TLSHandshake(info *tls.ClientHelloInfo) (ce *tls.Certificate, err error) {
	dial, err := tls.Dial("tcp", info.ServerName+":443", nil)
	if err != nil {
		return nil, err
	}

	if err = dial.Handshake(); err != nil {
		panic(err)
	}
	state := dial.ConnectionState()
	certificates := state.PeerCertificates
	temp := *certificates[0]
	return TLSCertificateForTLS(temp, CA)
}

func TLSCertificateForTLS(template x509.Certificate, cert *x509.Certificate) (ce *tls.Certificate, err error) {
	_, pkb, ceb, err := TLSCertificateFor(template, cert)
	if err != nil {
		return nil, err
	}
	keyPair, err := tls.X509KeyPair(ceb, pkb)
	return &keyPair, err
}
func TLSCertificateFor(template x509.Certificate, cert *x509.Certificate) (ce *x509.Certificate, pkb, ceb []byte, err error) {
	if cert == nil {
		cert = &template
	}
	var derBytes []byte
	derBytes, err = x509.CreateCertificate(rand.Reader, &template, cert, &PK.PublicKey, PK) //DER 格式
	if err != nil {
		return
	}
	ce, err = x509.ParseCertificate(derBytes)

	bufCe := bytes.NewBuffer([]byte{})
	err = pem.Encode(bufCe, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		return
	}
	bufPk := bytes.NewBuffer([]byte{})
	err = pem.Encode(bufPk, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(PK)})
	if err != nil {
		return
	}
	ceb = bufCe.Bytes()
	pkb = bufPk.Bytes()
	return
}
