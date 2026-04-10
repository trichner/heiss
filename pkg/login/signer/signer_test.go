package signer

import (
	"encoding/hex"
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestHmacSigner_Sign(t *testing.T) {
	signer := &HmacSigner{[]byte("s3cr3t")}

	data := []byte("Hello World!")
	expectedSignature := "b13404f85ffcdc8b6b452b2b60b2986cd81c852365628bb1a1ce4c51ca484a9d"

	signature, err := signer.Sign(data)

	assert.NoError(t, err)
	assert.Equal(t, expectedSignature, hex.EncodeToString(signature))
}

func TestHmacSigner_Verify_signAndVerify(t *testing.T) {
	signer := &HmacSigner{[]byte("s3cr3t")}

	data := []byte("Hello World!")

	signature, err := signer.Sign(data)
	assert.NoError(t, err)

	// when
	err = signer.Verify(data, signature)

	// then
	assert.NoError(t, err)
}

func TestHmacSigner_Verify_badSignature(t *testing.T) {
	signer := &HmacSigner{[]byte("s3cr3t")}

	data := []byte("Hello World!")

	badSignature, err := hex.DecodeString("badbeef85ffcdc8b6b452b2b60b2986cd81c852365628bb1a1ce4c51ca484a9d")
	assert.NoError(t, err)

	// when
	err = signer.Verify(data, badSignature)

	// then
	assert.ErrorIs(t, err, ErrBadSignature)
}
