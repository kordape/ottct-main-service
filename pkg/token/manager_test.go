package token

import (
	"fmt"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerateJWT(t *testing.T) {

	t.Run("generate token", func(t *testing.T) {
		secretKey := "s+76qrrI53tqP666lQla3t+pitXrAsYQQbEn/55JQKA="

		m, err := NewManager(secretKey, "foo")
		assert.NoError(t, err)

		tokenString, err := m.GenerateJWT("test")
		assert.NoError(t, err)
		assert.NotEmpty(t, tokenString)

		//verify signature
		claims := TokenClaims{}
		tkn, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(secretKey), nil
		})
		assert.NoError(t, err)
		assert.NotEmpty(t, tkn)
		assert.Equal(t, "test", claims.User)
		assert.Equal(t, "foo", claims.StandardClaims.Issuer)
	})

}
