package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func createJWT(claims map[string]any) (string, error) {
	privPem, err := os.ReadFile(os.Args[1])
	if err != nil {
		return "", err
	}

	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privPem)
	if err != nil {
		return "", err
	}

	return jwt.
		NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims(claims)).
		SignedString(privKey)
}

func validateJWT(signed string) (map[string]any, error) {
	pubPem, err := os.ReadFile(os.Args[2])
	if err != nil {
		return nil, err
	}

	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubPem)
	if err != nil {
		return nil, err
	}

	token, err := jwt.Parse(signed, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("validation failed")
}

func main() {
	signed, err := createJWT(map[string]any{
		"msg": "hello",
	})
	if err != nil {
		log.Fatal(err)
	}

	claims, err := validateJWT(signed)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Validated claims: ", claims)
}
