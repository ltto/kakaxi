package kakaxi

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log"
	"net"
	"net/http"
)

func OnTLS(accept net.Conn) (err error) {
	_, _ = accept.Write([]byte("HTTP/1.1 200 Connection Established\n\n"))
	config := tls.Config{GetCertificate: TLSHandshake}
	conn := tls.Server(accept, &config)
	if err = conn.Handshake(); err != nil {
		return err
	}

	request, err := http.ReadRequest(bufio.NewReader(conn))
	if err != nil {
		log.Printf("http.ReadRequest err:%v\n", err)
		return err
	}
	resp, err := ProxyHTTP(CopyRequest(request))
	if err != nil {
		return err
	}
	Writer(conn, resp)
	_ = conn.Close()
	return
}

func TLSHandshake(info *tls.ClientHelloInfo) (ce *tls.Certificate, err error) {
	switch info.ServerName {
	case "accounts.google.com":
		return nil, errors.New("不支持的域名:" + info.ServerName)
	}
	dial, err := tls.Dial("tcp", info.ServerName+":443", nil)
	if err != nil {
		log.Printf("tls.Dial(%s:443) err:%v\n", info.ServerName, err)
		return nil, err
	}

	if err = dial.Handshake(); err != nil {
		log.Printf("tls.Dial Handshake(%s:443) err:%v\n", info.ServerName, err)
		return nil, err
	}
	state := dial.ConnectionState()
	certificates := state.PeerCertificates
	temp := *certificates[0]
	return TLSCertificateForTLS(temp, CA)
}

func TLSCertificateForTLS(template x509.Certificate, cert *x509.Certificate) (ce *tls.Certificate, err error) {
	_, pkb, ceb, err := TLSCertificateFor(template, cert)
	if err != nil {
		log.Printf("TLSCertificateFor err:%v\n", err)
		return nil, err
	}
	keyPair, err := tls.X509KeyPair(ceb, pkb)
	if err != nil {
		log.Printf("tls.X509KeyPair err:%v\n", err)
		return nil, err
	}
	return &keyPair, err
}
func TLSCertificateFor(template x509.Certificate, cert *x509.Certificate) (ce *x509.Certificate, pkb, ceb []byte, err error) {
	if cert == nil {
		cert = &template
	}
	var derBytes []byte
	derBytes, err = x509.CreateCertificate(rand.Reader, &template, cert, &PK.PublicKey, PK) //DER 格式
	if err != nil {
		log.Printf("x509.CreateCertificate err:%v\n", err)
		return
	}
	ce, err = x509.ParseCertificate(derBytes)

	bufCe := bytes.NewBuffer([]byte{})
	err = pem.Encode(bufCe, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	if err != nil {
		log.Printf("pem.Encode err:%v\n", err)
		return
	}
	bufPk := bytes.NewBuffer([]byte{})
	err = pem.Encode(bufPk, &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(PK)})
	if err != nil {
		log.Printf("pem.Encode err:%v\n", err)
		return
	}
	ceb = bufCe.Bytes()
	pkb = bufPk.Bytes()
	return
}
