package mtls

import (
	"crypto/x509"
	"math/big"
	"time"
)

// CertificateOptions holds options for x509.CreateCertificate
type CertificateOptions struct {
	Organization          []string
	OrganizationalUnit    []string
	CommonName            string
	OCSPServer            []string
	IssuingCertificateURL []string
	SerialNumber          *big.Int
	NotAfter              time.Time
	NotBefore             time.Time
	SignatureAlgorithm    x509.SignatureAlgorithm
	PublicKeyAlgorithm    x509.PublicKeyAlgorithm
	ExtKeyUsage           []x509.ExtKeyUsage
	KeyUsage              x509.KeyUsage
	IsCA                  bool
}

// CertificateOrganizationalUnit set OrganizationalUnit in certificate subject
func CertificateOrganizationalUnit(s ...string) CertificateOption {
	return func(o *CertificateOptions) {
		o.OrganizationalUnit = s
	}
}

// CertificateOrganization set Organization in certificate subject
func CertificateOrganization(s ...string) CertificateOption {
	return func(o *CertificateOptions) {
		o.Organization = s
	}
}

// CertificateCommonName set CommonName in certificate subject
func CertificateCommonName(s string) CertificateOption {
	return func(o *CertificateOptions) {
		o.CommonName = s
	}
}

// CertificateOCSPServer set OCSPServer in certificate
func CertificateOCSPServer(s ...string) CertificateOption {
	return func(o *CertificateOptions) {
		o.OCSPServer = s
	}
}

// CertificateIssuingCertificateURL set IssuingCertificateURL in certificate
func CertificateIssuingCertificateURL(s ...string) CertificateOption {
	return func(o *CertificateOptions) {
		o.IssuingCertificateURL = s
	}
}

// CertificateSerialNumber set SerialNumber in certificate
func CertificateSerialNumber(n *big.Int) CertificateOption {
	return func(o *CertificateOptions) {
		o.SerialNumber = n
	}
}

// CertificateNotAfter set NotAfter in certificate
func CertificateNotAfter(t time.Time) CertificateOption {
	return func(o *CertificateOptions) {
		o.NotAfter = t
	}
}

// CertificateNotBefore set SerialNumber in certificate
func CertificateNotBefore(t time.Time) CertificateOption {
	return func(o *CertificateOptions) {
		o.NotBefore = t
	}
}

// CertificateExtKeyUsage set ExtKeyUsage in certificate
func CertificateExtKeyUsage(x ...x509.ExtKeyUsage) CertificateOption {
	return func(o *CertificateOptions) {
		o.ExtKeyUsage = x
	}
}

// CertificateSignatureAlgorithm set SignatureAlgorithm in certificate
func CertificateSignatureAlgorithm(alg x509.SignatureAlgorithm) CertificateOption {
	return func(o *CertificateOptions) {
		o.SignatureAlgorithm = alg
	}
}

// CertificatePublicKeyAlgorithm set PublicKeyAlgorithm in certificate
func CertificatePublicKeyAlgorithm(alg x509.PublicKeyAlgorithm) CertificateOption {
	return func(o *CertificateOptions) {
		o.PublicKeyAlgorithm = alg
	}
}

// CertificateKeyUsage set KeyUsage in certificate
func CertificateKeyUsage(u x509.KeyUsage) CertificateOption {
	return func(o *CertificateOptions) {
		o.KeyUsage = u
	}
}

// CertificateIsCA set IsCA in certificate
func CertificateIsCA(b bool) CertificateOption {
	return func(o *CertificateOptions) {
		o.IsCA = b
	}
}

// CertificateOption func signature
type CertificateOption func(*CertificateOptions)

func NewCertificateOptions(opts ...CertificateOption) CertificateOptions {
	options := CertificateOptions{}
	for _, o := range opts {
		o(&options)
	}
	if options.SerialNumber == nil {
		options.SerialNumber = big.NewInt(time.Now().UnixNano())
	}
	if options.NotBefore.IsZero() {
		options.NotBefore = time.Now()
	}
	if options.NotAfter.IsZero() {
		options.NotAfter = time.Now().Add(10 * time.Minute)
	}
	if options.SignatureAlgorithm == x509.UnknownSignatureAlgorithm {
		options.SignatureAlgorithm = x509.PureEd25519
	}
	if options.PublicKeyAlgorithm == x509.UnknownPublicKeyAlgorithm {
		options.PublicKeyAlgorithm = x509.Ed25519
	}
	if options.ExtKeyUsage == nil {
		options.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth}
		if options.IsCA {
			options.ExtKeyUsage = append(options.ExtKeyUsage, x509.ExtKeyUsageOCSPSigning, x509.ExtKeyUsageTimeStamping)
		}
	}

	if options.KeyUsage == 0 {
		options.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature
		if options.IsCA {
			options.KeyUsage = x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature | x509.KeyUsageDataEncipherment | x509.KeyUsageCertSign
		}
	}

	return options
}
