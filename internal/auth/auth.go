package auth

import (
	"errors"
	"fmt"
	"runtime"
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
		Issuer: "ShalomGate",
		IssuedAt: jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(expiresIn)), Subject: id.String(),
	}

	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, regClaims)
	token, err := tok.SignedString([]byte(secretKey))

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

//The auth token is sent as a cookie to the server, which is extracted by the middleware
func GetBearerToken(cookie, secret string) (string, error) {
	id, err := ValidateToken(cookie, secret)

	if err != nil {
		return "", fmt.Errorf("%w", err)
	}

	return id.String(), nil
}
