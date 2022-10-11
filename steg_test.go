package main

import (
	"os"
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

func TestConcealReveal(t *testing.T) {
	t.Log(mockInput(t, "hi there you"))
}

func mockInput(t *testing.T, str string) string {
	input := []byte(str)
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("creating pipe: %v", err)
	}
	defer w.Close()

	_, err = w.Write(input)
	if err != nil {
		t.Fatalf("writing input: %v", err)
	}

	defer func(v *os.File) { os.Stdin = v }(os.Stdin)
	os.Stdin = r

	return getInput("Enter public message")
}
