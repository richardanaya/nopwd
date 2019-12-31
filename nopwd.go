package nopwd

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/google/uuid"
)

type NoPwd struct {
	url    string
	secret string
	ttl    int64
}

func NewNoPwd(url, secret string, ttl int64) NoPwd {
	return NoPwd{
		url:    url,
		secret: secret,
		ttl:    ttl,
	}
}

func (self NoPwd) GenerateCodeLink(email string) (string, error) {
	code, err := self.generateJWT(email)
	if err != nil {
		return "", err
	}
	return self.url + "?code=" + code, nil
}

func (self NoPwd) generateJWT(email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti":   uuid.New(),
		"email": email,
		"iat":   time.Now().Unix(),
		"nbf":   time.Now().Unix(),
		"exp":   time.Now().Unix() + (self.ttl * 60),
		"iss":   self.url,
	})
	jwtToken, err := token.SignedString([]byte(self.secret))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (self NoPwd) ValidateCode(code string) (bool, string, error) {
	return self.validateCodeAtTime(code, time.Now().Unix())
}

func (self NoPwd) validateCodeAtTime(code string, currenTimeUnix int64) (bool, string, error) {
	token, err := jwt.Parse(code, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(self.secret), nil
	})
	if err != nil {
		return false, "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		issuer := claims["iss"].(string)
		expirationTime := int64(claims["exp"].(float64))
		notBeforeTime := int64(claims["nbf"].(float64))
		email := claims["email"].(string)
		if currenTimeUnix > expirationTime {
			return false, "", fmt.Errorf("token has expired")
		}
		if currenTimeUnix < notBeforeTime {
			return false, "", fmt.Errorf("token used before valid")
		}
		if issuer != self.url {
			return false, "", fmt.Errorf("token is not for this website")
		}
		return true, email, nil
	} else {
		return false, "", fmt.Errorf("token is not valid")
	}
}
