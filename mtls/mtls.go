package mtls // import "go.unistack.org/micro/v3/mtls"

import (
	"bytes"
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"sync"
)

var bp = newBPool()

type bpool struct {
	pool sync.Pool
}

func newBPool() *bpool {
	var bp bpool
	bp.pool.New = alloc
	return &bp
}

func alloc() interface{} {
	return &bytes.Buffer{}
}

func (bp *bpool) Get() *bytes.Buffer {
	return bp.pool.Get().(*bytes.Buffer)
}

func (bp *bpool) Put(buf *bytes.Buffer) {
	buf.Reset()
	bp.pool.Put(buf)
}

// NewCA creates new CA keypair
func NewCA(opts ...CertificateOption) ([]byte, crypto.PrivateKey, error) {
	options := NewCertificateOptions(opts...)

	crtreq := &x509.CertificateRequest{
		Subject: pkix.Name{
			Organization:       options.Organization,
			OrganizationalUnit: options.OrganizationalUnit,
			CommonName:         options.CommonName,
		},
		SignatureAlgorithm: options.SignatureAlgorithm,
	}

	pemcsr, pemkey, err := newCsr(crtreq)
	if err != nil {
		return nil, nil, err
	}

	pemcrt, err := SignCSR(pemcsr, nil, pemkey, opts...)
	if err != nil {
		return nil, nil, err
	}

	return pemcrt, pemkey, nil
}

func NewIntermediate(cacrt *x509.Certificate, cakey crypto.PrivateKey, opts ...CertificateOption) ([]byte, crypto.PrivateKey, error) {
	options := &CertificateOptions{}
	for _, o := range opts {
		o(options)
	}

	crtreq := &x509.CertificateRequest{
		Subject: pkix.Name{
			Organization:       options.Organization,
			OrganizationalUnit: options.OrganizationalUnit,
			CommonName:         options.CommonName,
		},
		SignatureAlgorithm: options.SignatureAlgorithm,
	}

	pemcsr, pemkey, err := newCsr(crtreq)
	if err != nil {
		return nil, nil, err
	}

	pemcrt, err := SignCSR(pemcsr, cacrt, cakey)
	if err != nil {
		return nil, nil, err
	}

	return pemcrt, pemkey, nil
}

// SignCSR sign certificate request and return signed pubkey
func SignCSR(rawcsr []byte, cacrt *x509.Certificate, cakey crypto.PrivateKey, opts ...CertificateOption) ([]byte, error) {
	if cacrt == nil {
		opts = append(opts, CertificateIsCA(false))
	}

	options := NewCertificateOptions(opts...)

	csr, err := x509.ParseCertificateRequest(rawcsr)
	if err == nil {
		err = csr.CheckSignature()
	}
	if err != nil {
		return nil, err
	}

	tpl := &x509.Certificate{
		Signature:             csr.Signature,
		SignatureAlgorithm:    csr.SignatureAlgorithm,
		PublicKeyAlgorithm:    csr.PublicKeyAlgorithm,
		PublicKey:             csr.PublicKey,
		SerialNumber:          options.SerialNumber,
		OCSPServer:            options.OCSPServer,
		IssuingCertificateURL: options.IssuingCertificateURL,
		Subject:               csr.Subject,
		NotBefore:             options.NotBefore,
		NotAfter:              options.NotAfter,
		KeyUsage:              options.KeyUsage,
		ExtKeyUsage:           options.ExtKeyUsage,
		BasicConstraintsValid: true,
		IsCA:                  options.IsCA,
	}

	if !options.IsCA {
		cacrt = tpl
	} else {
		tpl.Issuer = cacrt.Subject
	}

	crt, err := x509.CreateCertificate(rand.Reader, tpl, cacrt, csr.PublicKey, cakey)
	if err != nil {
		return nil, err
	}

	return crt, nil
}

// NewCertificateRequest create new certificate signing request and return key, csr in byte slice and err
func NewCertificateRequest(opts ...CertificateOption) ([]byte, crypto.PrivateKey, error) {
	options := NewCertificateOptions(opts...)

	crtreq := &x509.CertificateRequest{
		Subject: pkix.Name{
			Organization:       options.Organization,
			OrganizationalUnit: options.OrganizationalUnit,
			CommonName:         options.CommonName,
		},
		SignatureAlgorithm: options.SignatureAlgorithm,
	}

	return newCsr(crtreq)
}

// newCsr returns CSR and private key
func newCsr(crtreq *x509.CertificateRequest) ([]byte, crypto.PrivateKey, error) {
	_, key, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	csr, err := x509.CreateCertificateRequest(rand.Reader, crtreq, key)
	if err != nil {
		return nil, nil, err
	}
	return csr, key, nil
}

// ServerOptions holds server specific options
type ServerOptions struct {
	ServerName string
	RootCAs    []string
	ClientCAs  []string
}

// ServerOption func signature
type ServerOption func(*ServerOptions)

func NewServerConfig(src *tls.Config) *tls.Config {
	dst := src.Clone()
	dst.InsecureSkipVerify = true
	dst.MinVersion = tls.VersionTLS13
	dst.ClientAuth = tls.VerifyClientCertIfGiven
	return dst
}

func DecodeCrtKey(rawcrt []byte, rawkey []byte) (*x509.Certificate, crypto.PrivateKey, error) {
	var crt *x509.Certificate
	var key crypto.PrivateKey
	var err error

	crt, err = DecodeCrt(rawcrt)
	if err == nil {
		key, err = DecodeKey(rawkey)
	}

	if err != nil {
		return nil, nil, err
	}

	return crt, key, nil
}

func DecodeCrt(rawcrt []byte) (*x509.Certificate, error) {
	pemcrt, _ := pem.Decode(rawcrt)
	return x509.ParseCertificate(pemcrt.Bytes)
}

func EncodeCrt(crts ...*x509.Certificate) ([]byte, error) {
	var err error
	buf := bp.Get()
	defer bp.Put(buf)
	for _, crt := range crts {
		if err = pem.Encode(buf, &pem.Block{Type: "CERTIFICATE", Bytes: crt.Raw}); err != nil {
			return nil, err
		}
	}
	return buf.Bytes(), nil
}

func EncodeCsr(csr *x509.Certificate) ([]byte, error) {
	buf := bp.Get()
	defer bp.Put(buf)
	if err := pem.Encode(buf, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csr.Raw}); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeKey(rawkey []byte) (crypto.PrivateKey, error) {
	pemkey, _ := pem.Decode(rawkey)
	return x509.ParsePKCS8PrivateKey(pemkey.Bytes)
}

func EncodeKey(privkey crypto.PrivateKey) ([]byte, error) {
	buf := bp.Get()
	defer bp.Put(buf)
	enckey, err := x509.MarshalPKCS8PrivateKey(privkey)
	if err == nil {
		err = pem.Encode(buf, &pem.Block{Type: "PRIVATE KEY", Bytes: enckey})
	}
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
