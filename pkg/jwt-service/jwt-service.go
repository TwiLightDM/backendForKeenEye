package jwt_service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTService struct {
	Key         string
	AccessTime  time.Duration
	RefreshTime time.Duration
}

func NewJWTService(key string, accessTime, refreshTime time.Duration) *JWTService {
	return &JWTService{
		Key:         key,
		AccessTime:  accessTime,
		RefreshTime: refreshTime,
	}
}

func (s *JWTService) generateJWT(data map[string]any, expiration int64) (string, error) {
	claims := jwt.MapClaims{
		"exp": expiration,
	}

	for key, value := range data {
		if key == "id" {
			claims["sub"] = value
		} else {
			claims[key] = value
		}
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.Key))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

func (s *JWTService) GenerateRefreshJWT(data map[string]any) (string, error) {
	expiration := time.Now().Add(s.RefreshTime).Unix()
	delete(data, "exp")
	return s.generateJWT(data, expiration)

}

func (s *JWTService) GenerateAccessJWT(data map[string]any) (string, error) {
	expiration := time.Now().Add(s.AccessTime).Unix()
	delete(data, "exp")
	return s.generateJWT(data, expiration)
}

func (s *JWTService) ParseJWT(tokenString string) (map[string]any, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.Key), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", UnexpectedSigningMethodError)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return nil, fmt.Errorf("token expired: %w", LifetimeIsOverError)
		}
	}

	data := make(map[string]any)
	for key, value := range claims {
		data[key] = value
	}

	return data, nil
}
