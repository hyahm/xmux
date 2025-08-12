package xmux

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

// 生成证书， domain 只会读取第一个
func GenerateCertificate(certFile, keyFile string, domain ...string) error {
	dm := "localhost"
	if len(domain) > 0 {
		dm = domain[0]
	}
	addrs, err := net.LookupHost(dm)
	if err != nil {
		return err
	}
	iPAddresses := make([]net.IP, 0)
	for _, addr := range addrs {
		iPAddresses = append(iPAddresses, net.ParseIP(addr))
	}
	// 1. 生成私钥
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	// 2. 构造证书模板
	template := x509.Certificate{
		SerialNumber: big.NewInt(2024),
		Subject: pkix.Name{
			Organization: []string{dm},
			CommonName:   dm,
		},
		NotBefore:             time.Now().Add(-time.Hour), // 允许时钟偏差
		NotAfter:              time.Now().AddDate(10, 0, 0),
		BasicConstraintsValid: true,
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		DNSNames:              []string{dm},
		IPAddresses:           iPAddresses,
	}

	// 3. 自签名：用自己的私钥给自己签发
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return err
	}

	// 4. 编码并写入文件
	certOut, err := os.Create(certFile)
	if err != nil {
		return err
	}
	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certDER})
	certOut.Close()

	keyOut, err := os.Create(keyFile)
	if err != nil {
		return err
	}
	pem.Encode(keyOut, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(priv),
	})
	keyOut.Close()

	log.Printf("✅ 证书已生成：%s /%s\n", certFile, keyFile)
	return nil
}

func createTLS() {
	fi, err := os.Stat("keys")
	if err != nil {
		os.MkdirAll("keys", 0755)
	} else {
		if !fi.IsDir() {
			panic("exsit file")
		}
	}
	sn := time.Now().Unix()
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(sn),
		Subject: pkix.Name{
			Organization: []string{"scs"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 5},
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	privCa, _ := rsa.GenerateKey(rand.Reader, 1024)
	createCertificateFile("ca", ca, privCa, ca, nil)
	server := &x509.Certificate{
		SerialNumber: big.NewInt(sn),
		Subject: pkix.Name{
			Organization:       []string{"scs"},
			OrganizationalUnit: []string{"hyahm"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	hosts := make([]net.IP, 0)
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}
	for _, value := range addrs {
		if ipnet, ok := value.(*net.IPNet); ok {
			if ipnet.IP.To4() != nil {
				hosts = append(hosts, ipnet.IP)
			}
		}
	}

	server.IPAddresses = append(server.IPAddresses, hosts...)

	privSer, _ := rsa.GenerateKey(rand.Reader, 1024)
	createCertificateFile("server", server, privSer, ca, privCa)
	// client := &x509.Certificate{
	// 	SerialNumber: big.NewInt(sn),
	// 	Subject: pkix.Name{
	// 		Organization:       []string{"scs"},
	// 		OrganizationalUnit: []string{"hyahm"},
	// 	},
	// 	NotBefore:    time.Now(),
	// 	NotAfter:     time.Now().AddDate(10, 0, 0),
	// 	SubjectKeyId: []byte{1, 2, 3, 4, 7},
	// 	ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
	// 	KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	// }
	// privCli, _ := rsa.GenerateKey(rand.Reader, 1024)
	// CreateCertificateFile("client", client, privCli, ca, privCa)

}

func createCertificateFile(name string, cert *x509.Certificate, key *rsa.PrivateKey, caCert *x509.Certificate, caKey *rsa.PrivateKey) {
	name = filepath.Join("keys", name)
	priv := key
	pub := &priv.PublicKey
	privPm := priv
	if caKey != nil {
		privPm = caKey
	}
	ca_b, err := x509.CreateCertificate(rand.Reader, cert, caCert, pub, privPm)
	if err != nil {
		log.Println("create failed ", err)
		return
	}
	ca_f := name + ".pem"
	var certificate = &pem.Block{Type: "CERTIFICATE",
		Headers: map[string]string{},
		Bytes:   ca_b}
	ca_b64 := pem.EncodeToMemory(certificate)
	ioutil.WriteFile(ca_f, ca_b64, 0744)

	priv_f := name + ".key"
	priv_b := x509.MarshalPKCS1PrivateKey(priv)
	ioutil.WriteFile(priv_f, priv_b, 0744)
	var privateKey = &pem.Block{Type: "PRIVATE KEY",
		Headers: map[string]string{},
		Bytes:   priv_b}
	priv_b64 := pem.EncodeToMemory(privateKey)

	ioutil.WriteFile(priv_f, priv_b64, 0744)

}

func GenRsa(keyFile, certFile string) {
	// openssl genrsa -out private.pem 2048
	// openssl rsa -in private.pem -pubout -out public.pem
	// openssl pkeyutl -encrypt -inkey public.pem  -pubin -in f.txt -out fe.txt
	// openssl pkeyutl -decrypt -inkey private.pem -in fe.txt -out ffff.txt
	// 设置证书的有效期
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour) // 有效期为 365 天

	// 设置证书的主题
	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		log.Fatalf("Failed to generate serial number: %v", err)
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization:  []string{"My Company"},
			Country:       []string{"CN"},
			Province:      []string{"My Province"},
			Locality:      []string{"My City"},
			StreetAddress: []string{"My Street"},
			PostalCode:    []string{"123456"},
			CommonName:    "localhost", // 证书的通用名称，可以是域名或 IP 地址
		},
		NotBefore: notBefore,
		NotAfter:  notAfter,
		// 设置证书的用途
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		BasicConstraintsValid: true,
		IsCA:                  true, // 设置为自签名证书
	}

	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatalf("Failed to generate private key: %v", err)
	}

	// 使用私钥生成证书
	certificateBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		log.Fatalf("Failed to create certificate: %v", err)
	}

	// 保存证书到文件
	certOut, err := os.Create(certFile)
	if err != nil {
		log.Fatalf("Failed to open server.crt for writing: %v", err)
	}
	defer certOut.Close()

	// 使用 PEM 格式保存证书
	pem.Encode(certOut, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certificateBytes,
	})

	// 保存私钥到文件
	keyOut, err := os.OpenFile(keyFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Failed to open server.key for writing: %v", err)
	}
	defer keyOut.Close()

	// 使用 PEM 格式保存私钥
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	pem.Encode(keyOut, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	fmt.Println("Certificate and private key generated successfully!")
}

// RSA解密
func RsaDecryptFromBase64(s string, priviteKeyPath string) ([]byte, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	file, err := os.Open(priviteKeyPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//获取文件内容
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//对密文进行解密
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, b)
}

// RSA解密
func RsaDecryptFromString(s string, priviteKeyPath string) ([]byte, error) {
	b := []byte(s)
	file, err := os.Open(priviteKeyPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//获取文件内容
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//对密文进行解密
	return rsa.DecryptPKCS1v15(rand.Reader, privateKey, b)
}

// RSA加密
func RSA_Encrypt(plainText []byte, path string) ([]byte, error) {
	//打开文件
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	//读取文件的内容
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	//x509解码

	publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//类型断言
	publicKey := publicKeyInterface.(*rsa.PublicKey)
	//对明文进行加密
	return rsa.EncryptPKCS1v15(rand.Reader, publicKey, plainText)
}
