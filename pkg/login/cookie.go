package login

import (
	"encoding/base64"
	"fmt"
	"strings"
)

type SignedToken struct {
	Payload   []byte
	Signature []byte
}

func CreateCookie(s *SignedToken) string {
	payload := encodeBase64(s.Payload)
	signature := encodeBase64(s.Signature)
	token := payload + "." + signature
	return token
}

func ParseCookie(s string) (SignedToken, error) {
	splits := strings.SplitN(s, ".", 2)
	if len(splits) != 2 {
		return SignedToken{}, fmt.Errorf("invalid format: %q", s)
	}

	payload, err := parseBase64(splits[0])
	if err != nil {
		return SignedToken{}, fmt.Errorf("cannot parse payload of %q: %w", splits[0], err)
	}

	signature, err := parseBase64(splits[1])
	if err != nil {
		return SignedToken{}, fmt.Errorf("cannot parse signature of %q: %w", splits[1], err)
	}

	return SignedToken{Payload: payload, Signature: signature}, nil
}

func parseBase64(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

func encodeBase64(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}
