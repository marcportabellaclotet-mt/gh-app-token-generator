package helpers

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func LookupEnvString(key string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return ""
}

func CheckEnvVars(key string, envVar string) string {
	if key != "" {
		return key
	} else {
		return LookupEnvString(envVar)
	}
}

func LoadPEMFromBytes(key []byte) (*rsa.PrivateKey, error) {
	b, _ := pem.Decode(key)
	if b != nil {
		key = b.Bytes
	}
	parsedKey, err := x509.ParsePKCS1PrivateKey(key)
	if err != nil {
		return nil, fmt.Errorf("private key should be a PKCS1 key; parse error: %v", err)
	}
	return parsedKey, nil
}

func IssueJWTFromPEM(key *rsa.PrivateKey, appID string) string {

	claims := &jwt.StandardClaims{
		IssuedAt:  time.Now().Add(-1 * time.Minute).Unix(),
		ExpiresAt: time.Now().Add(time.Minute * 10).Unix(),
		Issuer:    appID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	ss, err := token.SignedString(key)
	if err != nil {
		fmt.Println(err)
		return ""
	}

	return ss
}

func String(v string) *string {
	return &v
}
