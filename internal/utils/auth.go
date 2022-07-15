package utils

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nickklius/go-loyalty/config"
)

type Token struct {
	AccessToken        string
	AccessTokenExpires int64
}

func CreateToken(userID string, cfg config.Auth) (*Token, error) {
	t := &Token{
		AccessTokenExpires: time.Now().Add(time.Minute * time.Duration(cfg.AccessLifeTime)).Unix(),
	}

	accessClaims := jwt.MapClaims{
		"exp":     t.AccessTokenExpires,
		"user_id": userID,
	}

	accessWithClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err := accessWithClaims.SignedString([]byte(cfg.Secret))
	if err != nil {
		return nil, err
	}

	t.AccessToken = accessToken

	return t, nil
}

func ValidateToken(r *http.Request, cfg *config.Auth) (*jwt.Token, error) {
	tokenString := ExtractToken(r)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	if _, ok := token.Claims.(jwt.MapClaims); !ok && !token.Valid {
		return nil, errors.New("expired token")
	}

	return token, nil
}

func ExtractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")

	strArr := strings.Split(bearerToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}
