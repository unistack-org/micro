package id_test

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	id "github.com/unistack-org/micro/v3/util/id"
)

func TestGenerate(t *testing.T) {
	t.Run("short alphabet", func(t *testing.T) {
		alphabet := ""
		_, err := id.New(id.Alphabet(alphabet), id.Size(32))
		assert.Error(t, err, "should return error if the alphabet is too small")
	})

	t.Run("long alphabet", func(t *testing.T) {
		alphabet := strings.Repeat("a", 256)
		_, err := id.New(id.Alphabet(alphabet), id.Size(32))
		assert.Error(t, err, "should return error if the alphabet is too long")
	})

	t.Run("negative ID length", func(t *testing.T) {
		_, err := id.New(id.Alphabet("abcdef"), id.Size(-1))
		assert.Error(t, err, "should return error if the requested ID length is invalid")
	})

	t.Run("happy path", func(t *testing.T) {
		alphabet := "abcdef"
		id, err := id.New(id.Alphabet(alphabet), id.Size(6))
		assert.NoError(t, err, "shouldn't return error")
		assert.Len(t, id, 6, "should return ID of requested length")
		for _, r := range id {
			assert.True(t, strings.ContainsRune(alphabet, r), "should use given alphabet")
		}
	})

	t.Run("works with unicode", func(t *testing.T) {
		alphabet := "🚀💩🦄🤖"
		id, err := id.New(id.Alphabet(alphabet), id.Size(6))
		assert.NoError(t, err, "shouldn't return error")
		assert.Equal(t, utf8.RuneCountInString(id), 6, "should return ID of requested length")
		for _, r := range id {
			assert.True(t, strings.ContainsRune(alphabet, r), "should use given alphabet")
		}
	})
}

func TestNew(t *testing.T) {
	t.Run("negative ID length", func(t *testing.T) {
		_, err := id.New(id.Size(-1))
		assert.Error(t, err, "should return error if the requested ID length is invalid")
	})

	t.Run("happy path", func(t *testing.T) {
		nid, err := id.New()
		assert.NoError(t, err, "shouldn't return error")
		assert.Len(t, nid, id.DefaultSize, "should return ID of default length")
	})

	t.Run("custom length", func(t *testing.T) {
		id, err := id.New(id.Size(6))
		assert.NoError(t, err, "shouldn't return error")
		assert.Len(t, id, 6, "should return ID of requested length")
	})
}
