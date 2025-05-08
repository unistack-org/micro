package id

import (
	"crypto/rand"
	"errors"

	"github.com/google/uuid"
	nanoid "github.com/matoous/go-nanoid"
)

var generatedNode [6]byte

func init() {
	if _, err := rand.Read(generatedNode[:]); err != nil {
		panic(err)
	}
}

type Type int

const (
	TypeUnspecified Type = iota
	TypeNanoid
	TypeUUIDv7
	TypeUUIDv8
)

// DefaultNanoidAlphabet is the alphabet used for ID characters by default
var DefaultNanoidAlphabet = "6789BCDFGHJKLMNPQRTWbcdfghjkmnpqrtwz"

// DefaultNanoidSize is the size used for ID by default
// To get uuid like collision specify 21
var DefaultNanoidSize = 16

type Generator struct {
	opts Options
}

func (g *Generator) MustNew() string {
	id, err := g.New()
	if err != nil {
		panic(err)
	}
	return id
}

func (g *Generator) New() (string, error) {
	switch g.opts.Type {
	case TypeNanoid:
		if len(g.opts.NanoidAlphabet) == 0 || len(g.opts.NanoidAlphabet) > 255 {
			return "", errors.New("invalid option, NanoidAlphabet must not be empty and contain no more than 255 chars")
		}
		if g.opts.NanoidSize <= 0 {
			return "", errors.New("invalid option, NanoidSize must be positive integer")
		}

		return nanoid.Generate(g.opts.NanoidAlphabet, g.opts.NanoidSize)
	case TypeUUIDv7:
		uid, err := uuid.NewV7()
		if err != nil {
			return "", err
		}
		return uid.String(), nil
	case TypeUUIDv8:
		return "", errors.New("unsupported uuid version v8")
	}
	return "", errors.New("invalid option, Type unspecified")
}

// New returns new id or error
func New(opts ...Option) (string, error) {
	options := NewOptions(opts...)

	switch options.Type {
	case TypeNanoid:
		if len(options.NanoidAlphabet) == 0 || len(options.NanoidAlphabet) > 255 {
			return "", errors.New("invalid option, NanoidAlphabet must not be empty and contain no more than 255 chars")
		}
		if options.NanoidSize <= 0 {
			return "", errors.New("invalid option, NanoidSize must be positive integer")
		}
		return nanoid.Generate(options.NanoidAlphabet, options.NanoidSize)
	case TypeUUIDv7:
		uid, err := uuid.NewV7()
		if err != nil {
			return "", err
		}
		return uid.String(), nil
	case TypeUUIDv8:
		return "", errors.New("unsupported uuid version v8")
	}

	return "", errors.New("invalid option, Type unspecified")
}

func ToUUID(s string) uuid.UUID {
	return uuid.MustParse(s)
}

// Must is the same as New but fatals on error
func MustNew(opts ...Option) string {
	id, err := New(opts...)
	if err != nil {
		panic(err)
	}
	return id
}

// Options contains id deneration options
type Options struct {
	Type           Type
	NanoidAlphabet string
	NanoidSize     int
	UUIDNode       [6]byte
}

// Option func signature
type Option func(*Options)

// WithNanoidAlphabet specifies alphabet to use
func WithNanoidAlphabet(alphabet string) Option {
	return func(o *Options) {
		o.NanoidAlphabet = alphabet
	}
}

// WithNanoidSize specifies generated id size
func WithNanoidSize(size int) Option {
	return func(o *Options) {
		o.NanoidSize = size
	}
}

// WithUUIDNode specifies node component for UUIDv8
func WithUUIDNode(node [6]byte) Option {
	return func(o *Options) {
		o.UUIDNode = node
	}
}

// NewOptions returns new Options struct filled by opts
func NewOptions(opts ...Option) Options {
	options := Options{
		Type:           TypeUUIDv7,
		NanoidAlphabet: DefaultNanoidAlphabet,
		NanoidSize:     DefaultNanoidSize,
		UUIDNode:       generatedNode,
	}

	for _, o := range opts {
		o(&options)
	}

	return options
}
