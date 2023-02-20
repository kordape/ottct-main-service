package token

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

const defaultTokenExpiry = 30 * time.Minute

type Manager struct {
	secretHS256Key string
	expiry         time.Duration
	issuer         string
}

type Option func(*Manager)

func WithExpiry(expiry time.Duration) Option {
	return func(m *Manager) {
		m.expiry = expiry
	}
}

func NewManager(secretKey string, issuer string) (*Manager, error) {
	m := Manager{
		secretHS256Key: secretKey,
		expiry:         defaultTokenExpiry,
		issuer:         issuer,
	}

	if err := m.validate(); err != nil {
		return &m, err
	}

	return &m, nil
}

func (m *Manager) validate() error {
	if m.secretHS256Key == "" {
		return fmt.Errorf("Token manager validation error: invalid secret key")
	}

	if m.issuer == "" {
		return fmt.Errorf("Token manager validation error: invalid issuer")
	}

	return nil
}

// TokenClaims - Claims for a JWT access token.
type TokenClaims struct {
	User string `json:"user"`
	jwt.StandardClaims
}

func (m *Manager) GenerateJWT(user string) (string, error) {
	// Set-up claims
	claims := TokenClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(m.expiry).Unix(),
			Issuer:    m.issuer,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(m.secretHS256Key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
