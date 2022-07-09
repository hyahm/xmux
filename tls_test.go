package xmux

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"os"
	"testing"
)

func TestRsa(t *testing.T) {
	GenRsa("pri.key", "pub.key", "ca.crt")
}

func TestDecrypt(t *testing.T) {
	s := "OCT+5EHncL1sZg0BRniPz5x3A0MTv24pltkjzyoYE3ucf46MtCZGNs+NojDD+6Rl7z1wvP6jJeAlmaUYZukuOfvf7sdcwA4KZWCZycc25RIavoFzNrPbjqPpsf0u2k/BCleP7Dc8d3Ca34ZQlrlIIGxwj13pBKFHSMrR8YZtgvn9EueXM4/ENpq2MY9eamm2X5e6IsmfFl8J+L33U4NNhokXCBo/5wcXfwFiIBzSJXf9R+NRWjCITVxgyctiWyx4jtu8NrKYaiNDUhuyrsf1jmzgPgAKEOeC7VuwV0HRqa9UgNsAV7P6K8ZBMsvrXoEMWB6HQRZ28aLhtevrbFeymA=="
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		t.Fatal(err)
	}
	file, err := os.Open("pri.pem")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	//获取文件内容
	info, _ := file.Stat()
	buf := make([]byte, info.Size())
	file.Read(buf)
	//pem解码
	block, _ := pem.Decode(buf)
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		panic(err)
	}
	//对密文进行解密
	plainText, _ := rsa.DecryptPKCS1v15(rand.Reader, privateKey, b)
	t.Log(string(plainText))

}
