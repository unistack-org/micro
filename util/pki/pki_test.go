package pki

import (
	"crypto/rand"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"testing"
	"time"
)

func TestPrivateKey(t *testing.T) {
	_, _, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
}

func TestCA(t *testing.T) {
	pub, priv, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	serialNumberMax := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberMax)
	if err != nil {
		t.Fatal(err)
	}

	cert, key, err := CA(
		KeyPair(pub, priv),
		Subject(pkix.Name{
			Organization: []string{"test"},
		}),
		DNSNames("localhost"),
		IPAddresses(net.ParseIP("127.0.0.1")),
		SerialNumber(serialNumber),
		NotBefore(time.Now().Add(time.Minute*-1)),
		NotAfter(time.Now().Add(time.Minute)),
	)
	if err != nil {
		t.Fatal(err)
	}

	asn1Key, _ := pem.Decode(key)
	if asn1Key == nil {
		t.Fatal(err)
	}
	if asn1Key.Type != "PRIVATE KEY" {
		t.Fatal("invalid key type")
	}
	decodedKey, err := x509.ParsePKCS8PrivateKey(asn1Key.Bytes)
	if err != nil {
		t.Fatal(err)
	} else if decodedKey == nil {
		t.Fatal("empty key")
	}

	asn1Cert, _ := pem.Decode(cert)
	if asn1Cert == nil {
		t.Fatal(err)
	}

	/*
		pool := x509.NewCertPool()

		x509cert, err := x509.ParseCertificate(asn1Cert.Bytes)
		if err != nil {
			t.Fatal(err)
		}


		chains, err := x509cert.Verify(x509.VerifyOptions{
			Roots: pool,
		})
		if err != nil {
			t.Fatal(err)
		}

		if len(chains) != 1 {
			t.Fatal("CA should have 1 cert in chain")
		}
	*/
}

func TestCSR(t *testing.T) {
	pub, priv, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}

	csr, err := CSR(
		Subject(
			pkix.Name{
				CommonName:         "testnode",
				Organization:       []string{"microtest"},
				OrganizationalUnit: []string{"super-testers"},
			},
		),
		DNSNames("localhost"),
		IPAddresses(net.ParseIP("127.0.0.1")),
		KeyPair(pub, priv),
	)
	if err != nil {
		t.Fatal(err)
	}

	asn1csr, _ := pem.Decode(csr)
	if asn1csr == nil {
		t.Fatal(err)
	}

	decodedcsr, err := x509.ParseCertificateRequest(asn1csr.Bytes)
	if err != nil {
		t.Fatal(err)
	}

	expected := pkix.Name{
		CommonName:         "testnode",
		Organization:       []string{"microtest"},
		OrganizationalUnit: []string{"super-testers"},
	}
	if decodedcsr.Subject.String() != expected.String() {
		t.Fatalf("%s != %s", decodedcsr.Subject.String(), expected.String())
	}
}
