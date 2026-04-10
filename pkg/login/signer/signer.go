package signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"errors"
	"fmt"
)

var ErrBadSignature = errors.New("bad signature")

func NewHmacSigner(key []byte) *HmacSigner {
	return &HmacSigner{key: key}
}

type HmacSigner struct {
	key []byte
}

func (h *HmacSigner) Sign(data []byte) ([]byte, error) {
	signature, err := h.digest(data)
	if err != nil {
		return nil, fmt.Errorf("cannot sign: %w", err)
	}
	return signature, nil
}

func (h *HmacSigner) Verify(data, signature []byte) error {
	expectedSignature, err := h.digest(data)
	if err != nil {
		return fmt.Errorf("cannot verify signature: %w", err)
	}

	if !hmac.Equal(expectedSignature, signature) {
		return ErrBadSignature
	}
	return nil
}

func (h *HmacSigner) digest(data []byte) ([]byte, error) {
	signer := hmac.New(sha256.New, h.key)

	_, err := signer.Write(data)
	if err != nil {
		return nil, err
	}

	return signer.Sum(nil), nil
}
