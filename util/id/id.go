package id // import "go.unistack.org/micro/v4/util/id"

import (
	"context"
	"crypto/rand"
	"errors"
	"math"

	"go.unistack.org/micro/v4/logger"
)

// DefaultAlphabet is the alphabet used for ID characters by default
var DefaultAlphabet = []rune("6789BCDFGHJKLMNPQRTWbcdfghjkmnpqrtwz")

// DefaultSize is the size used for ID by default
// To get uuid like collision specify 21
var DefaultSize = 16

// getMask generates bit mask used to obtain bits from the random bytes that are used to get index of random character
// from the alphabet. Example: if the alphabet has 6 = (110)_2 characters it is sufficient to use mask 7 = (111)_2
func getMask(alphabetSize int) int {
	for i := 1; i <= 8; i++ {
		mask := (2 << uint(i)) - 1
		if mask >= alphabetSize-1 {
			return mask
		}
	}
	return 0
}

// New returns new id or error
func New(opts ...Option) (string, error) {
	options := NewOptions(opts...)

	if len(options.Alphabet) == 0 || len(options.Alphabet) > 255 {
		return "", errors.New("alphabet must not be empty and contain no more than 255 chars")
	}
	if options.Size <= 0 {
		return "", errors.New("size must be positive integer")
	}

	chars := options.Alphabet

	mask := getMask(len(chars))
	// estimate how many random bytes we will need for the ID, we might actually need more but this is tradeoff
	// between average case and worst case
	ceilArg := 1.6 * float64(mask*options.Size) / float64(len(options.Alphabet))
	step := int(math.Ceil(ceilArg))

	id := make([]rune, options.Size)
	bytes := make([]byte, step)
	for j := 0; ; {
		_, err := rand.Read(bytes)
		if err != nil {
			return "", err
		}
		for i := 0; i < step; i++ {
			currByte := bytes[i] & byte(mask)
			if currByte < byte(len(chars)) {
				id[j] = chars[currByte]
				j++
				if j == options.Size {
					return string(id[:options.Size]), nil
				}
			}
		}
	}
}

// Must is the same as New but fatals on error
func Must(opts ...Option) string {
	id, err := New(opts...)
	if err != nil {
		logger.Fatal(context.TODO(), err)
	}
	return id
}

// Options contains id deneration options
type Options struct {
	Alphabet []rune
	Size     int
}

// Option func signature
type Option func(*Options)

// Alphabet specifies alphabet to use
func Alphabet(alphabet string) Option {
	return func(o *Options) {
		o.Alphabet = []rune(alphabet)
	}
}

// Size specifies id size
func Size(size int) Option {
	return func(o *Options) {
		o.Size = size
	}
}

// NewOptions returns new Options struct filled by opts
func NewOptions(opts ...Option) Options {
	options := Options{
		Alphabet: DefaultAlphabet,
		Size:     DefaultSize,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
