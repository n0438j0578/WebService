package model

import (
	"errors"

	"github.com/dgrijalva/jwt-go"
)

var (
	mssp   = []byte("m2gkz-+n5z1bXgHwhnMf4qxct8@gn^&;")
	Issuer = "Time-Recorder"
)

type Claims struct {
	Username string
	Uid      string
	Role     string
	Token    string
	jwt.StandardClaims
}

func DecodeSession(session string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(session, claims, func(token *jwt.Token) (interface{}, error) {
		return mssp, nil
	})
	if err != nil {
		return claims, err
	}
	if claims.Issuer != Issuer {
		return claims, errors.New("Fail!!")
	}
	return claims, nil
}

func NewSessionTime(id, username string, role string, token string, expires int64) string {
	claims := Claims{
		Username: username,
		Uid:      id,
		Role:     role,
		Token:    token,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires,
			Issuer:    Issuer,
		},
	}
	jwttoken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := jwttoken.SignedString(mssp)
	if err != nil {
		return ""
	}
	return ss
}
