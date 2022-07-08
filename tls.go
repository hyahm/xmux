package xmux

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/base64"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"time"
)

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

	for _, ip := range hosts {
		server.IPAddresses = append(server.IPAddresses, ip)
	}

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

func GenRsa(prikey, pubkey, crtkey string) error {
	// openssl genrsa -out private.pem 2048
	// openssl rsa -in private.pem -pubout -out public.pem
	// openssl pkeyutl -encrypt -inkey public.pem  -pubin -in f.txt -out fe.txt
	// openssl pkeyutl -decrypt -inkey private.pem -in fe.txt -out ffff.txt
	pri, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}

	derText := x509.MarshalPKCS1PrivateKey(pri)
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derText,
	}
	f, err := os.OpenFile(prikey, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	pem.Encode(f, block)

	// 生成私钥

	ptext, err := x509.MarshalPKIXPublicKey(&pri.PublicKey)
	if err != nil {
		return err
	}
	pblock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: ptext,
	}

	f, err = os.OpenFile(pubkey, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	pem.Encode(f, pblock)
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject:      pkix.Name{CommonName: "xmux"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		KeyUsage:     x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDigitalSignature,
	}
	if ip := net.ParseIP("localhost"); ip != nil {
		template.IPAddresses = append(template.IPAddresses, ip)
	} else {
		template.DNSNames = append(template.DNSNames, "localhost")
	}
	cert, err := x509.CreateCertificate(rand.Reader, &template, &template, &pri.PublicKey, pri)
	if err != nil {
		return err
	}
	cblock := &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert,
	}
	f, err = os.OpenFile(crtkey, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return pem.Encode(f, cblock)
}

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

//RSA加密
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
