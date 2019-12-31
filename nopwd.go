package nopwd

import (
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/google/uuid"
)

type NoPwd struct {
	secret string
}

func NewNoPwd(secret string) NoPwd {
	return NoPwd{
		secret: secret,
	}
}

func (self NoPwd) GenerateLoginLink(url, email string, ttl int64) (string, error) {
	code, err := self.GenerateLoginCode(email, ttl)
	if err != nil {
		return "", err
	}
	return url + "?login_code=" + code, nil
}

func (self NoPwd) GenerateLoginCode(email string, ttl int64) (string, error) {
	code, err := self.generateJWT(email, "login", ttl)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (self NoPwd) GenerateAPICode(email string, ttl int64) (string, error) {
	code, err := self.generateJWT(email, "api", ttl)
	if err != nil {
		return "", err
	}
	return code, nil
}

func (self NoPwd) generateJWT(email, code_type string, ttl int64) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"jti":       uuid.New(),
		"email":     email,
		"iat":       time.Now().Unix(),
		"nbf":       time.Now().Unix(),
		"exp":       time.Now().Unix() + (ttl * 60),
		"code_type": code_type,
	})
	jwtToken, err := token.SignedString([]byte(self.secret))
	if err != nil {
		return "", err
	}

	return jwtToken, nil
}

func (self NoPwd) ValidateLoginCode(code string) (bool, string, error) {
	return self.validateCodeAtTime(code, "login", time.Now().Unix())
}

func (self NoPwd) ValidateAPICode(code string) (bool, string, error) {
	return self.validateCodeAtTime(code, "api", time.Now().Unix())
}

func (self NoPwd) validateCodeAtTime(code, codeType string, currenTimeUnix int64) (bool, string, error) {
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
		expirationTime := int64(claims["exp"].(float64))
		notBeforeTime := int64(claims["nbf"].(float64))
		email := claims["email"].(string)
		jwtCodeType := claims["code_type"].(string)
		println(jwtCodeType)
		println(codeType)
		if codeType != jwtCodeType {
			return false, "", fmt.Errorf("token has unexpected code type")
		}
		if currenTimeUnix > expirationTime {
			return false, "", fmt.Errorf("token has expired")
		}
		if currenTimeUnix < notBeforeTime {
			return false, "", fmt.Errorf("token used before valid")
		}
		return true, email, nil
	} else {
		return false, "", fmt.Errorf("token is not valid")
	}
}
