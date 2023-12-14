package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserLogin
type Claims struct {
	jwt.RegisteredClaims
	UserLogin string
	ClientID  string
}

type ClientInfo struct {
	UserLogin string
	ClientID  string
}

const tokenExp = time.Hour * 3
const secretKey = "supersecretkey"

// BuildJWTString создаёт токен и возвращает его в виде строки.
func BuildJWTString(login string, clientID string) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenExp)),
		},
		// собственное утверждение
		UserLogin: login,
		ClientID:  clientID,
	})

	// создаём строку токена
	// token.Method = &jwt.SigningMethodHMAC{}
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func GetClientInfo(tokenString string) (*ClientInfo, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(secretKey), nil
		})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	return &ClientInfo{
			UserLogin: claims.UserLogin,
			ClientID:  claims.ClientID,
		},
		nil
}

func GenerateUUID() string {
	return uuid.New().String()
}
