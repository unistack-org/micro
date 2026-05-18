package mtls

import (
	"crypto/ed25519"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"testing"
	"time"
)

func TestNewCa(t *testing.T) {
	bcrt, key, err := NewCA(
		CertificateOrganization("test_org"),
		CertificateOrganizationalUnit("test_unit"),
		CertificateIsCA(true),
	)
	if err != nil {
		t.Fatal(err)
	}

	if _, ok := key.(ed25519.PrivateKey); !ok {
		t.Fatalf("key is not ed25519")
	}

	crt, err := x509.ParseCertificate(bcrt)
	if err != nil {
		t.Fatal(err)
	}
	if !crt.IsCA {
		t.Fatalf("crt IsCA invalid %v", crt)
	}
	if crt.Subject.Organization[0] != "test_org" {
		t.Fatalf("crt subject invalid %v", crt.Subject)
	}
	if crt.Subject.OrganizationalUnit[0] != "test_unit" {
		t.Fatalf("crt subject invalid %v", crt.Subject)
	}
}

func TestNewIntermediate(t *testing.T) {
	bcrt, cakey, err := NewCA(
		CertificateOrganization("test_org"),
		CertificateOrganizationalUnit("test_unit"),
	)
	if err != nil {
		t.Fatal(err)
	}
	cacrt, err := x509.ParseCertificate(bcrt)
	if err != nil {
		t.Fatal(err)
	}

	bcrt, ikey, err := NewIntermediate(cacrt, cakey,
		CertificateOrganization("test_org"),
		CertificateOrganizationalUnit("test_unit"),
	)
	if err != nil {
		t.Fatal(err)
	}
	_ = ikey
	icrt, err := x509.ParseCertificate(bcrt)
	if err != nil {
		t.Fatal(err)
	}

	if icrt.IsCA {
		t.Fatalf("crt IsCA invalid %v", icrt)
	}
	if icrt.Subject.Organization[0] != "test_org" {
		t.Fatalf("crt subject invalid %v", icrt.Subject)
	}
	if icrt.Subject.OrganizationalUnit[0] != "test_unit" {
		t.Fatalf("crt subject invalid %v", icrt.Subject)
	}
}

func TestNewCertificateRequest(t *testing.T) {
	csr, key, err := NewCertificateRequest(
		CertificateCommonName("test.local"),
		CertificateOrganization("myorg"),
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(csr) == 0 {
		t.Fatal("expected non-empty CSR")
	}
	if key == nil {
		t.Fatal("expected non-nil key")
	}
}

func TestNewServerConfig(t *testing.T) {
	src := &tls.Config{ServerName: "test"}
	dst := NewServerConfig(src)
	if dst == nil {
		t.Fatal("expected non-nil config")
	}
	if dst.MinVersion != tls.VersionTLS13 {
		t.Fatalf("unexpected MinVersion: %d", dst.MinVersion)
	}
	if dst.ClientAuth != tls.VerifyClientCertIfGiven {
		t.Fatalf("unexpected ClientAuth: %v", dst.ClientAuth)
	}
}

func TestEncodeDecodeKey(t *testing.T) {
	_, key, err := NewCA(CertificateIsCA(true))
	if err != nil {
		t.Fatal(err)
	}

	// EncodeKey
	pemkey, err := EncodeKey(key)
	if err != nil {
		t.Fatal(err)
	}
	if len(pemkey) == 0 {
		t.Fatal("expected non-empty PEM key")
	}

	// DecodeKey
	decodedKey, err := DecodeKey(pemkey)
	if err != nil {
		t.Fatal(err)
	}
	if decodedKey == nil {
		t.Fatal("expected non-nil decoded key")
	}
}

func TestEncodeCrt(t *testing.T) {
	bcrt, key, err := NewCA(CertificateIsCA(true))
	if err != nil {
		t.Fatal(err)
	}
	_ = key

	// bcrt is raw DER — parse it first
	crt, err := x509.ParseCertificate(bcrt)
	if err != nil {
		t.Fatal(err)
	}

	// EncodeCrt wraps crt.Raw in PEM
	pemcrt, err := EncodeCrt(crt)
	if err != nil {
		t.Fatal(err)
	}
	if len(pemcrt) == 0 {
		t.Fatal("expected non-empty PEM cert")
	}
}

func TestDecodeCrtAndDecodeCrtKey(t *testing.T) {
	bcrt, key, err := NewCA(CertificateIsCA(true))
	if err != nil {
		t.Fatal(err)
	}

	// bcrt is raw DER; PEM-encode it directly (not via EncodeCrt which uses crt.Raw)
	pemcrt := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: bcrt})

	pemkey, err := EncodeKey(key)
	if err != nil {
		t.Fatal(err)
	}

	// DecodeCrt
	decoded, err := DecodeCrt(pemcrt)
	if err != nil {
		t.Fatal(err)
	}
	if decoded == nil {
		t.Fatal("expected non-nil decoded cert")
	}

	// DecodeCrtKey
	crt2, key2, err := DecodeCrtKey(pemcrt, pemkey)
	if err != nil {
		t.Fatal(err)
	}
	if crt2 == nil || key2 == nil {
		t.Fatal("expected non-nil decoded cert/key pair")
	}
}

func TestEncodeCsr(t *testing.T) {
	bcrt, _, err := NewCA(CertificateIsCA(true))
	if err != nil {
		t.Fatal(err)
	}
	crt, err := x509.ParseCertificate(bcrt)
	if err != nil {
		t.Fatal(err)
	}
	pemcsr, err := EncodeCsr(crt)
	if err != nil {
		t.Fatal(err)
	}
	if len(pemcsr) == 0 {
		t.Fatal("expected non-empty PEM CSR")
	}
}

func TestSignCSR(t *testing.T) {
	// Create a CA
	bcacrt, cakey, err := NewCA(CertificateIsCA(true))
	if err != nil {
		t.Fatal(err)
	}
	cacrt, err := x509.ParseCertificate(bcacrt)
	if err != nil {
		t.Fatal(err)
	}

	// Create a CSR
	rawcsr, _, err := NewCertificateRequest(CertificateCommonName("leaf.local"))
	if err != nil {
		t.Fatal(err)
	}

	// Sign the CSR
	signedCrt, err := SignCSR(rawcsr, cacrt, cakey)
	if err != nil {
		t.Fatal(err)
	}
	if len(signedCrt) == 0 {
		t.Fatal("expected non-empty signed cert")
	}
}

func TestCertificateOptions(t *testing.T) {
	now := time.Now()
	serial := big.NewInt(42)

	opts := NewCertificateOptions(
		CertificateCommonName("test.local"),
		CertificateOrganization("org1"),
		CertificateOrganizationalUnit("unit1"),
		CertificateOCSPServer("http://ocsp.example.com"),
		CertificateIssuingCertificateURL("http://ica.example.com"),
		CertificateSerialNumber(serial),
		CertificateNotBefore(now),
		CertificateNotAfter(now.Add(time.Hour)),
		CertificateExtKeyUsage(x509.ExtKeyUsageServerAuth),
		CertificateSignatureAlgorithm(x509.PureEd25519),
		CertificatePublicKeyAlgorithm(x509.Ed25519),
		CertificateKeyUsage(x509.KeyUsageDigitalSignature),
	)

	if opts.CommonName != "test.local" {
		t.Errorf("CommonName = %q", opts.CommonName)
	}
	if opts.Organization[0] != "org1" {
		t.Errorf("Organization = %v", opts.Organization)
	}
	if opts.OrganizationalUnit[0] != "unit1" {
		t.Errorf("OrganizationalUnit = %v", opts.OrganizationalUnit)
	}
	if opts.OCSPServer[0] != "http://ocsp.example.com" {
		t.Errorf("OCSPServer = %v", opts.OCSPServer)
	}
	if opts.IssuingCertificateURL[0] != "http://ica.example.com" {
		t.Errorf("IssuingCertificateURL = %v", opts.IssuingCertificateURL)
	}
	if opts.SerialNumber.Cmp(serial) != 0 {
		t.Errorf("SerialNumber = %v", opts.SerialNumber)
	}
}
