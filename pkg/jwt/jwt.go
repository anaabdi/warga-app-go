package jwt

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	defaultIssuer = "warga-app-go"
)

func CreateToken(ttl time.Duration, claims map[string]interface{}, privateKey string) (string, error) {
	decodedPrivateKey, err := base64.StdEncoding.DecodeString(privateKey)
	if err != nil {
		return "", fmt.Errorf("could not decode key: %w", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(decodedPrivateKey)

	if err != nil {
		return "", fmt.Errorf("create: parse key: %w", err)
	}

	now := time.Now().UTC()

	mapClaims := make(jwt.MapClaims)
	mapClaims["sub"] = claims["sub"]
	mapClaims["iss"] = defaultIssuer
	mapClaims["exp"] = now.Add(ttl).Unix()
	mapClaims["iat"] = now.Unix()
	mapClaims["nbf"] = now.Unix()
	mapClaims["acc_type"] = claims["acc_type"]

	if v, found := claims["nxtrid"]; found {
		mapClaims["nxtrid"] = v
	}

	if v, found := claims["sm"]; found {
		mapClaims["sm"] = v
	}

	if v, found := claims["scope"]; found {
		mapClaims["scope"] = v
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, mapClaims).SignedString(key)

	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

func ValidateToken(token string, publicKey string) (jwt.MapClaims, error) {
	decodedPublicKey, err := base64.StdEncoding.DecodeString(publicKey)
	if err != nil {
		return nil, fmt.Errorf("could not decode: %w", err)
	}

	key, err := jwt.ParseRSAPublicKeyFromPEM(decodedPublicKey)

	if err != nil {
		return nil, fmt.Errorf("validate: parse key: %w", err)
	}

	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected method: %s", t.Header["alg"])
		}
		return key, nil
	})

	if err != nil {
		return nil, fmt.Errorf("validate: %w", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, fmt.Errorf("validate: invalid token")
	}

	return claims, nil
}
