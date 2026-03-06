package auth

import (
	"time"
	"github.com/golang-jwt/jwt"
)

type MyCustomClaims struct {
    UserID string `json:"user_id"`
    Role   string `json:"role"`
    jwt.StandardClaims
}

func GenerateToken(signature []byte, userId string, role string) (string, error) {
    claims := MyCustomClaims{
        UserID: userId,
        Role:   role,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(15 * time.Minute).Unix(),
            Issuer:    "apartment_sys",
            IssuedAt:  time.Now().Unix(),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(signature)
}

func ValidateToken(signature []byte, tokenString string) (*MyCustomClaims, error) {
		token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, jwt.ErrSignatureInvalid
				}
				return signature, nil
		})

		if err != nil {
				return nil, err
		}

		if claims, ok := token.Claims.(*MyCustomClaims); ok && token.Valid {
				return claims, nil
		}
		return nil, jwt.ErrInvalidKey
}