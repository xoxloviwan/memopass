package jwt

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const Bearer = "Bearer "
const TokenExp = time.Hour * 24
const SecretKey = "supersecretkey"

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

// BuildJWT создаёт токен и возвращает его в виде строки.
func BuildJWT(login string, id int) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
			Subject:   login,
		},
		UserID: id,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

func GetUser(tokenString string) (*int, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(SecretKey), nil
		})
	if err != nil {
		return nil, err
	}

	if !token.Valid || claims.Subject == "" || claims.UserID == 0 {
		return nil, errors.New("invalid token")
	}

	return &claims.UserID, nil
}
