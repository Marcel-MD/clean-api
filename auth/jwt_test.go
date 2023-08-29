package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func TestGenerateJWT(t *testing.T) {
	userId := "123"
	roles := []string{"admin", "user"}
	lifespan := time.Hour
	secret := "mysecret"

	token, err := GenerateAccessToken(userId, roles, lifespan, secret)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if !parsedToken.Valid {
		t.Errorf("generated token is not valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Errorf("unexpected token claims type")
	}

	if !claims["authorized"].(bool) {
		t.Errorf("token authorization claim is incorrect")
	}

	if claims["user_id"].(string) != userId {
		t.Errorf("token user_id claim is incorrect")
	}

	if len(claims["roles"].([]interface{})) != len(roles) {
		t.Errorf("token roles claim is incorrect")
	}
}

func TestExtractId(t *testing.T) {
	userId := "123"
	roles := []string{"admin", "user"}
	lifespan := time.Hour
	secret := "mysecret"

	token, err := GenerateAccessToken(userId, roles, lifespan, secret)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	parsedToken, err := Validate(token, secret)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	extractedId, err := ExtractId(parsedToken)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if extractedId != userId {
		t.Errorf("extracted user_id is incorrect")
	}
}

func TestExtractIdAndRoles(t *testing.T) {
	userId := "123"
	roles := []string{"admin", "user"}
	lifespan := time.Hour
	secret := "mysecret"

	token, err := GenerateAccessToken(userId, roles, lifespan, secret)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	parsedToken, err := Validate(token, secret)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	extractedId, extractedRoles, err := ExtractIdAndRoles(parsedToken)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if extractedId != userId {
		t.Errorf("extracted user_id is incorrect")
	}

	if len(extractedRoles) != len(roles) {
		t.Errorf("extracted roles are incorrect")
	}

	for i, v := range roles {
		if extractedRoles[i] != v {
			t.Errorf("extracted roles are incorrect")
		}
	}
}
