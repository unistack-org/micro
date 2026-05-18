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

func TestSign(t *testing.T) {
	// Generate CA cert and key
	pub, priv, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	serialMax := new(big.Int).Lsh(big.NewInt(1), 128)
	serial, err := rand.Int(rand.Reader, serialMax)
	if err != nil {
		t.Fatal(err)
	}
	caCert, caKey, err := CA(
		KeyPair(pub, priv),
		Subject(pkix.Name{Organization: []string{"test-ca"}}),
		SerialNumber(serial),
		NotBefore(time.Now().Add(-time.Minute)),
		NotAfter(time.Now().Add(time.Hour)),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Generate a leaf CSR
	leafPub, leafPriv, err := GenerateKey()
	if err != nil {
		t.Fatal(err)
	}
	csrPEM, err := CSR(
		KeyPair(leafPub, leafPriv),
		Subject(pkix.Name{CommonName: "leaf.local"}),
		DNSNames("leaf.local"),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Sign the CSR with the CA
	leafSerial, err := rand.Int(rand.Reader, serialMax)
	if err != nil {
		t.Fatal(err)
	}
	leafCert, err := Sign(caCert, caKey, csrPEM,
		SerialNumber(leafSerial),
		NotBefore(time.Now().Add(-time.Minute)),
		NotAfter(time.Now().Add(time.Hour)),
	)
	if err != nil {
		t.Fatal(err)
	}

	asn1Cert, _ := pem.Decode(leafCert)
	if asn1Cert == nil {
		t.Fatal("expected valid PEM cert from Sign")
	}
	decodedCert, err := x509.ParseCertificate(asn1Cert.Bytes)
	if err != nil {
		t.Fatal(err)
	}
	if decodedCert.Subject.CommonName != "leaf.local" {
		t.Errorf("unexpected CN: %s", decodedCert.Subject.CommonName)
	}
}

func TestSignInvalidCA(t *testing.T) {
	_, err := Sign([]byte("not-pem"), []byte("not-pem"), []byte("not-pem"))
	if err == nil {
		t.Fatal("expected error for invalid CA cert PEM")
	}
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
