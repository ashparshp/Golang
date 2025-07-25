package urlsigner

import (
	"fmt"
	"strings"

	goalone "github.com/bwmarrin/go-alone"
)

type Signer struct {
	Secret []byte
}

func (s *Signer) GenerateTokenFromString(data string) string {
	var urlToSign string

	crypt := goalone.New(s.Secret, goalone.Timestamp)

	if strings.Contains(data, "?") {
		urlToSign = fmt.Sprintf("%s&hash=", data)
	} else {
		urlToSign = fmt.Sprintf("%s?hash=", data)
	}

	tokenBytes := crypt.Sign([]byte(urlToSign))
	token := string(tokenBytes)
	return token
}

func (s *Signer) VerifyToken(token string) bool {
	// Implementation for verifying a token
	// This is a placeholder; actual implementation will depend on the signing algorithm used
	return token == "generated_token"
}

func (s *Signer) Expired(token string, minutesUntillExpire int) bool {
	// Implementation for checking if a token is expired
	// This is a placeholder; actual implementation will depend on the token structure
	return false
}
