package mtls

import (
	"crypto/ed25519"
	"crypto/x509"
	"testing"
)

func TestNewCa(t *testing.T) {
	bcrt, key, err := NewCA(
		CertificateOrganization("test_org"),
		CertificateOrganizationalUnit("test_unit"),
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
	if crt.IsCA {
		t.Fatalf("crt IsCA invalid %v", crt)
	}
	if crt.Subject.Organization[0] != "test_org" {
		t.Fatalf("crt subject invalid %v", crt.Subject)
	}
	if crt.Subject.OrganizationalUnit[0] != "test_unit" {
		t.Fatalf("crt subject invalid %v", crt.Subject)
	}
}
