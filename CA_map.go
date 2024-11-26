package kakaxi

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"
)

// CAMap 存储签名算法对应的证书密钥对
var CAMap = map[x509.SignatureAlgorithm]CertKeyPair{}

// DomainCertCache 存储域名证书缓存
var DomainCertCache struct {
	sync.RWMutex
	certs map[string]DomainCert
}

// DomainCert 域名证书结构
type DomainCert struct {
	Cert      *x509.Certificate
	CertPEM   []byte
	ExpiresAt time.Time
}

func init() {
	CAMap = generateAllCertKeyPairs()
	DomainCertCache.certs = make(map[string]DomainCert)
}

type CertKeyPair struct {
	Certificate *x509.Certificate
	PrivateKey  interface{}
	PublicKey   interface{} // 添加公钥字段
	CertPEM     []byte
	KeyPEM      []byte
	PubKeyPEM   []byte // 可选：PEM格式的公钥
}

func generateAllCertKeyPairs() map[x509.SignatureAlgorithm]CertKeyPair {
	result := make(map[x509.SignatureAlgorithm]CertKeyPair)

	algorithms := []struct {
		sigAlg     x509.SignatureAlgorithm
		genKeyFunc func() (interface{}, interface{}, error)
		desc       string
	}{
		// 与之前相同的算法列表...
		{x509.SHA1WithRSA, rsaKeyGen(2048), "SHA1WithRSA (legacy)"},
		{x509.SHA256WithRSA, rsaKeyGen(2048), "SHA256WithRSA"},
		{x509.SHA384WithRSA, rsaKeyGen(3072), "SHA384WithRSA"},
		{x509.SHA512WithRSA, rsaKeyGen(4096), "SHA512WithRSA"},
		{x509.ECDSAWithSHA1, ecdsaKeyGen(elliptic.P256()), "ECDSAWithSHA1 (legacy)"},
		{x509.ECDSAWithSHA256, ecdsaKeyGen(elliptic.P256()), "ECDSAWithSHA256"},
		{x509.ECDSAWithSHA384, ecdsaKeyGen(elliptic.P384()), "ECDSAWithSHA384"},
		{x509.ECDSAWithSHA512, ecdsaKeyGen(elliptic.P521()), "ECDSAWithSHA512"},
		{x509.SHA256WithRSAPSS, rsaKeyGen(2048), "SHA256WithRSAPSS"},
		{x509.SHA384WithRSAPSS, rsaKeyGen(3072), "SHA384WithRSAPSS"},
		{x509.SHA512WithRSAPSS, rsaKeyGen(4096), "SHA512WithRSAPSS"},
		{x509.PureEd25519, ed25519KeyGen(), "PureEd25519"},
	}

	for _, alg := range algorithms {
		privateKey, publicKey, err := alg.genKeyFunc()
		if err != nil {
			log.Printf("生成密钥对失败 %v: %v", alg.desc, err)
			continue
		}

		template := &x509.Certificate{
			SerialNumber: big.NewInt(1),
			Subject: pkix.Name{
				Organization: []string{"Test Org"},
				CommonName:   fmt.Sprintf("Test %v Cert", alg.desc),
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().AddDate(1, 0, 0),
			KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
			BasicConstraintsValid: true,
			SignatureAlgorithm:    alg.sigAlg,
		}

		derBytes, err := x509.CreateCertificate(rand.Reader, template, template, publicKey, privateKey)
		if err != nil {
			log.Printf("创建证书失败 %v: %v", alg.desc, err)
			continue
		}

		cert, err := x509.ParseCertificate(derBytes)
		if err != nil {
			log.Printf("解析证书失败 %v: %v", alg.desc, err)
			continue
		}

		certPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: derBytes,
		})

		// 生成密钥的PEM编码
		var keyPEM, pubKeyPEM []byte
		switch k := privateKey.(type) {
		case *rsa.PrivateKey:
			keyPEM = pem.EncodeToMemory(&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(k),
			})
			pubKeyPEM = pem.EncodeToMemory(&pem.Block{
				Type:  "RSA PUBLIC KEY",
				Bytes: x509.MarshalPKCS1PublicKey(&k.PublicKey),
			})
		case *ecdsa.PrivateKey:
			keyBytes, err := x509.MarshalECPrivateKey(k)
			if err != nil {
				log.Printf("编码ECDSA私钥失败 %v: %v", alg.desc, err)
				continue
			}
			keyPEM = pem.EncodeToMemory(&pem.Block{
				Type:  "EC PRIVATE KEY",
				Bytes: keyBytes,
			})
			pubKeyBytes, err := x509.MarshalPKIXPublicKey(&k.PublicKey)
			if err != nil {
				log.Printf("编码ECDSA公钥失败 %v: %v", alg.desc, err)
				continue
			}
			pubKeyPEM = pem.EncodeToMemory(&pem.Block{
				Type:  "PUBLIC KEY",
				Bytes: pubKeyBytes,
			})
		case ed25519.PrivateKey:
			keyBytes, err := x509.MarshalPKCS8PrivateKey(k)
			if err != nil {
				log.Printf("编码Ed25519私钥失败 %v: %v", alg.desc, err)
				continue
			}
			keyPEM = pem.EncodeToMemory(&pem.Block{
				Type:  "PRIVATE KEY",
				Bytes: keyBytes,
			})
			pubKeyBytes, err := x509.MarshalPKIXPublicKey(k.Public())
			if err != nil {
				log.Printf("编码Ed25519公钥失败 %v: %v", alg.desc, err)
				continue
			}
			pubKeyPEM = pem.EncodeToMemory(&pem.Block{
				Type:  "PUBLIC KEY",
				Bytes: pubKeyBytes,
			})
		}

		result[alg.sigAlg] = CertKeyPair{
			Certificate: cert,
			PrivateKey:  privateKey,
			PublicKey:   publicKey,
			CertPEM:     certPEM,
			KeyPEM:      keyPEM,
			PubKeyPEM:   pubKeyPEM,
		}
	}

	return result
}

// 辅助函数：生成RSA密钥对
func rsaKeyGen(bits int) func() (interface{}, interface{}, error) {
	return func() (interface{}, interface{}, error) {
		key, err := rsa.GenerateKey(rand.Reader, bits)
		if err != nil {
			return nil, nil, err
		}
		return key, &key.PublicKey, nil
	}
}

// 辅助函数：生成ECDSA密钥对
func ecdsaKeyGen(curve elliptic.Curve) func() (interface{}, interface{}, error) {
	return func() (interface{}, interface{}, error) {
		key, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		return key, &key.PublicKey, nil
	}
}

// 辅助函数：生成Ed25519密钥对
func ed25519KeyGen() func() (interface{}, interface{}, error) {
	return func() (interface{}, interface{}, error) {
		pub, priv, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return nil, nil, err
		}
		return priv, pub, nil
	}
}

// 辅助函数 - 根据签名算法获取对应的密钥对
func GetKeyPairForAlgorithm(sigAlg x509.SignatureAlgorithm) (*CertKeyPair, error) {
	pair, exists := CAMap[sigAlg]
	if !exists {
		return nil, fmt.Errorf("不支持的签名算法: %v", sigAlg)
	}
	return &pair, nil
}

// 示例：如何使用该函数创建证书
func CreateCertificateWithAlgorithm(template *x509.Certificate, parent *x509.Certificate) (*x509.Certificate, error) {
	// 获取对应算法的密钥对
	keyPair, err := GetKeyPairForAlgorithm(template.SignatureAlgorithm)
	if err != nil {
		return nil, err
	}

	// 如果是自签名证书
	if parent == nil {
		parent = template
	}

	// 创建证书
	derBytes, err := x509.CreateCertificate(rand.Reader, template, parent, keyPair.PublicKey, keyPair.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("创建证书失败: %v", err)
	}

	// 解析并返回证书
	return x509.ParseCertificate(derBytes)
}

// GetOrCreateDomainCert 获取或创建域名证书
func GetOrCreateDomainCert(domain string, sigAlg x509.SignatureAlgorithm) (*DomainCert, error) {
	// 先尝试从缓存读取
	DomainCertCache.RLock()
	if cert, exists := DomainCertCache.certs[domain]; exists {
		DomainCertCache.RUnlock()
		// 检查证书是否过期
		if time.Now().Before(cert.ExpiresAt) {
			return &cert, nil
		}
	} else {
		DomainCertCache.RUnlock()
	}

	// 没有找到或已过期，需要创建新证书
	DomainCertCache.Lock()
	defer DomainCertCache.Unlock()

	// 双重检查，避免并发创建
	if cert, exists := DomainCertCache.certs[domain]; exists {
		if time.Now().Before(cert.ExpiresAt) {
			return &cert, nil
		}
	}

	// 获取CA证书和密钥
	keyPair, err := GetKeyPairForAlgorithm(sigAlg)
	if err != nil {
		return nil, err
	}

	// 创建域名证书模板
	template := &x509.Certificate{
		SerialNumber: big.NewInt(time.Now().UnixNano()),
		Subject: pkix.Name{
			CommonName: domain,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0), // 1年有效期
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		SignatureAlgorithm:    sigAlg,
		DNSNames:              []string{domain},
	}

	// 使用CA证书签发域名证书
	derBytes, err := x509.CreateCertificate(rand.Reader, template, keyPair.Certificate, keyPair.PublicKey, keyPair.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("创建域名证书失败: %v", err)
	}

	// 解析证书
	cert, err := x509.ParseCertificate(derBytes)
	if err != nil {
		return nil, fmt.Errorf("解析域名证书失败: %v", err)
	}

	// 生成PEM格式证书
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: derBytes,
	})

	// 创建新的域名证书记录
	domainCert := DomainCert{
		Cert:      cert,
		CertPEM:   certPEM,
		ExpiresAt: template.NotAfter,
	}

	// 更新缓存
	DomainCertCache.certs[domain] = domainCert

	return &domainCert, nil
}

// CleanExpiredCerts 清理过期的域名证书
func CleanExpiredCerts() {
	DomainCertCache.Lock()
	defer DomainCertCache.Unlock()

	now := time.Now()
	for domain, cert := range DomainCertCache.certs {
		if now.After(cert.ExpiresAt) {
			delete(DomainCertCache.certs, domain)
		}
	}
}

// GetCachedDomainsCount 获取缓存的域名证书数量
func GetCachedDomainsCount() int {
	DomainCertCache.RLock()
	defer DomainCertCache.RUnlock()
	return len(DomainCertCache.certs)
}

// GetCachedDomains 获取所有缓存的域名列表
func GetCachedDomains() []string {
	DomainCertCache.RLock()
	defer DomainCertCache.RUnlock()

	domains := make([]string, 0, len(DomainCertCache.certs))
	for domain := range DomainCertCache.certs {
		domains = append(domains, domain)
	}
	return domains
}

// 其他原有代码保持不变...
