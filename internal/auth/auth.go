package auth

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var params = &argon2id.Params{
	Memory:      128 * 1024,
	Iterations:  4,
	Parallelism: uint8(runtime.NumCPU()),
	SaltLength:  16,
	KeyLength:   16,
}

func CreateHash(password string) (string, error) {
	return argon2id.CreateHash(password, params)
}


func CheckHash(password, hash string) (bool, error) {
	return argon2id.ComparePasswordAndHash(password, hash)
}

/**
 * Used to make authentication/access tokens
*/
func MakeTokens(secretKey string, id uuid.UUID, expiresIn time.Duration) (string, error) {

	now := time.Now()
	regClaims := jwt.RegisteredClaims{
		Issuer: "ShalomGate", IssuedAt: jwt.NewNumericDate(now), ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)), Subject: id.String(),
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, regClaims)
	token, err := tok.SignedString(secretKey)

	if err != nil {
		return "", err
	}
	return token, nil
}


func ValidateToken(tokenString, secretKey string) (uuid.UUID, error) {
	claims := jwt.MapClaims{}

	// Parse the token using your secret
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})

	sub, ok := claims["sub"].(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("invalid subject claim")
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return uuid.Nil, err
	}

	return userID, nil
}

func GetBearerToken(headers http.Header, secret string) (string, error) {
	// extract the authorization header
	rawHeader := headers.Get("Authorization")
	if rawHeader == "" {
		return "", fmt.Errorf("authorization header is missing")
	}

	// split on spacec
	parts := strings.Fields(rawHeader)
	if len(parts) < 2 {
		return "", fmt.Errorf("authorization header is malformed (expected: Bearer <token>)")
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return "", fmt.Errorf("authorization type must be Bearer")
	}

	tokenString := parts[1]
	if tokenString == "" {
		return "", fmt.Errorf("bearer token is empty")
	}

	id, err := ValidateToken(tokenString, secret)
	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return id.String(), nil
}
