package xmux

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
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
