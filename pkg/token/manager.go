package token

import (
	"fmt"
	"strings"
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
	User uint `json:"user"`
	jwt.StandardClaims
}

func (m *Manager) GenerateJWT(userId uint) (string, error) {
	// Set-up claims
	claims := TokenClaims{
		User: userId,
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

func (m *Manager) VerifyJWT(tokenString string) error {
	claims := TokenClaims{}
	_, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.secretHS256Key), nil
	})

	if err != nil {
		return fmt.Errorf("Invalid token: %w", err)
	}

	return nil
}

func (m *Manager) GetClaimsFromJWT(auth string) (*TokenClaims, error) {
	token := strings.Split(auth, "Bearer ")[1]
	claims := TokenClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(m.secretHS256Key), nil
	})
	if err != nil {
		return nil, fmt.Errorf("error parsing token with claims: %w", err)
	}

	return &claims, nil
}
