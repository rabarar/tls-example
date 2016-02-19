package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/asn1"
	"encoding/pem"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"time"
)

type PEMType string

const (
	PEM_CERTIFICATE PEMType = "CERTIFICATE"
	PEM_PRIVATE_KEY PEMType = "RSA PRIVATE KEY"
	PEM_PUBLICE_KEY PEMType = "RSA PUBLIC KEY"

	PUB_PERMS  = 0666
	PRIV_PERMS = 0400

	PEM_EXTENSION = ".pem"
	RAW_EXTENSION = ".raw"
)

func main() {

	// Config the CA Cert
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
		IPAddresses:           []net.IP{[]byte{0, 0, 0, 0}, []byte{127, 0, 0, 1}},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 5},
		BasicConstraintsValid: true,
		IsCA:        true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}

	// Generate a keypair
	priv, _ := rsa.GenerateKey(rand.Reader, 1024)
	pub := &priv.PublicKey

	// Generate a Cert from the CA config and the keypair
	CACert, err := x509.CreateCertificate(rand.Reader, ca, ca, pub, priv)

	if err != nil {
		log.Println("create ca failed", err)
		return
	}

	// write the CA cert in PEM format
	CAFile := "ca_cert"
	log.Println("write to", CAFile)
	ioutil.WriteFile(CAFile+PEM_EXTENSION, byteToPEM(CACert, PEM_CERTIFICATE), PUB_PERMS)
	ioutil.WriteFile(CAFile+RAW_EXTENSION, CACert, PUB_PERMS)

	// write the CA private key in PEM format
	CAPrivFile := "ca_key"
	priv_b := x509.MarshalPKCS1PrivateKey(priv)
	log.Println("write to", CAPrivFile)
	ioutil.WriteFile(CAPrivFile+PEM_EXTENSION, byteToPEM(priv_b, PEM_PRIVATE_KEY), PRIV_PERMS)
	ioutil.WriteFile(CAPrivFile+RAW_EXTENSION, priv_b, PRIV_PERMS)

	// Config a Cert from the CA
	clientCert := &x509.Certificate{
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

	// Generate a keypair for the client Cert
	clientKey, _ := rsa.GenerateKey(rand.Reader, 1024)
	clientPubKey := &clientKey.PublicKey
	clientCertRaw, err2 := x509.CreateCertificate(rand.Reader, clientCert, ca, clientPubKey, priv)
	if err2 != nil {
		log.Println("create cert2 failed", err2)
		return
	}

	// Write the client Cert to PEM file
	clientCertFile := "client_cert"
	log.Println("write to", clientCertFile)
	ioutil.WriteFile(clientCertFile+PEM_EXTENSION, byteToPEM(clientCertRaw, PEM_CERTIFICATE), PUB_PERMS)
	ioutil.WriteFile(clientCertFile+RAW_EXTENSION, clientCertRaw, PUB_PERMS)

	// Write the client private Key to PEM file
	clientPrivKeyFile := "client_key"
	ClientPrivateKeyRaw := x509.MarshalPKCS1PrivateKey(clientKey)
	log.Println("write to", clientPrivKeyFile)
	ioutil.WriteFile(clientPrivKeyFile+PEM_EXTENSION, byteToPEM(ClientPrivateKeyRaw, PEM_PRIVATE_KEY), PRIV_PERMS)
	ioutil.WriteFile(clientPrivKeyFile+RAW_EXTENSION, ClientPrivateKeyRaw, PRIV_PERMS)

	// verify the Certs
	caCheck, _ := x509.ParseCertificate(CACert)
	clientCertCheck, _ := x509.ParseCertificate(clientCertRaw)

	err3 := clientCertCheck.CheckSignatureFrom(caCheck)
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

func byteToPEM(b []byte, t PEMType) []byte {
	var blk pem.Block = pem.Block{
		Type:    string(t),
		Headers: nil,
	}

	blk.Bytes = b
	data := pem.EncodeToMemory(&blk)
	return data
}
