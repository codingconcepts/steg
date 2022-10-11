package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	pt := "the quick brown fox jumps over the lazy dog"
	key := []byte("password")

	ct := encryptData(pt, key)
	act := decryptData(ct, key)

	assert.Equal(t, act, pt)
}

func TestEncode(t *testing.T) {
	pub := "hi"
	pri := "there"

	act := encode(pub, pri)
	assert.Equal(t, "h\u2060\u2060\u2060\u200c\u2060\u200c\u200c\u200b\u2060\u2060\u200c\u2060\u200c\u200c\u200c\u200b\u2060\u2060\u200c\u200c\u2060\u200c\u2060\u200b\u2060\u2060\u2060\u200c\u200c\u2060\u200c\u200b\u2060\u2060\u200c\u200c\u2060\u200c\u2060i", act)
}

func TestDecode(t *testing.T) {
	pub := "h\u2060\u2060\u2060\u200c\u2060\u200c\u200c\u200b\u2060\u2060\u200c\u2060\u200c\u200c\u200c\u200b\u2060\u2060\u200c\u200c\u2060\u200c\u2060\u200b\u2060\u2060\u2060\u200c\u200c\u2060\u200c\u200b\u2060\u2060\u200c\u200c\u2060\u200c\u2060i"

	act := decode(pub)
	assert.Equal(t, "there", act)
}
