package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"time"
)

func main() {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Country:            []string{"Hobbitsville"},
			Organization:       []string{"Bilbo"},
			OrganizationalUnit: []string{"Babbins"},
			Locality:           []string{"ca-locality"},
			Province:           []string{"ca-province"},
			StreetAddress:      []string{"ca-streetaddr"},
			PostalCode:         []string{"ca-postalcode"},
			SerialNumber:       "ca-serial",
			CommonName:         "ca-commonName",
			Names:              addASN1Names(1, 3, 6, 1),
			ExtraNames:         addASN1Names(1, 3, 6, 2),
		},
		//DNSNames:              []string{"0.0.0.0"},
		//IPAddresses:           []net.IP{[]byte{127, 0, 0, 1}},
		IPAddresses:           []net.IP{[]byte{0, 0, 0, 0}, []byte{127, 0, 0, 1}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 5},
		BasicConstraintsValid: true,
		IsCA:        true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	priv, _ := rsa.GenerateKey(rand.Reader, 1024)
	pub := &priv.PublicKey
	ca_b, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)
	if err != nil {
		log.Println("create ca failed", err)
		return
	}
	ca_f := "ca.pem"
	log.Println("write to", ca_f)
	ioutil.WriteFile(ca_f, ca_b, 0777)

	priv_f := "ca.key"
	priv_b := x509.MarshalPKCS1PrivateKey(priv)
	log.Println("write to", priv_f)
	ioutil.WriteFile(priv_f, priv_b, 0777)

	cert2 := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Country:            []string{"Narnia"},
			Organization:       []string{"Whizard"},
			OrganizationalUnit: []string{"Druids"},
			Locality:           []string{"client-locality"},
			Province:           []string{"client-province"},
			StreetAddress:      []string{"client-streetaddr"},
			PostalCode:         []string{"client-postalcode"},
			SerialNumber:       "client-serial",
			CommonName:         "client-commonName",
			//Names      []AttributeTypeAndValue
			//ExtraNames []AttributeTypeAndValue
		},
		DNSNames:     []string{"0.0.0.0"},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	priv2, _ := rsa.GenerateKey(rand.Reader, 1024)
	pub2 := &priv2.PublicKey
	cert2_b, err2 := x509.CreateCertificate(rand.Reader, cert2, ca, pub2, priv)
	if err2 != nil {
		log.Println("create cert2 failed", err2)
		return
	}

	cert2_f := "cert2.pem"
	log.Println("write to", cert2_f)
	ioutil.WriteFile(cert2_f, cert2_b, 0777)

	priv2_f := "cert2.key"
	priv2_b := x509.MarshalPKCS1PrivateKey(priv2)
	log.Println("write to", priv2_f)
	ioutil.WriteFile(priv2_f, priv2_b, 0777)

	ca_c, _ := x509.ParseCertificate(ca_b)
	cert2_c, _ := x509.ParseCertificate(cert2_b)

	err3 := cert2_c.CheckSignatureFrom(ca_c)
	log.Println("check signature", err3 == nil)
}

func addASN1Names(x ...int) []pkix.AttributeTypeAndValue {
	var asn1Vals pkix.AttributeTypeAndValue
	var daytime interface{}

	if true {
		daytime = "hi there ASN1"
	} else {
		daytime = time.Now()
		// Ignore return network errors.
	}
	mdata, _ := asn1.Marshal(daytime)
	asn1Vals.Type = x
	asn1Vals.Value = mdata

	return []pkix.AttributeTypeAndValue{asn1Vals}
}
