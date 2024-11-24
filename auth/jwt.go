package auth

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"time"
)

type JWTService struct {
	issuer            string
	signingMethod     jwt.SigningMethod
	signingMethodAlg  []string
	expiryTimeMinutes int
	signingKey        []byte
	parser            *jwt.Parser
}

type JWTClaims struct {
	jwt.RegisteredClaims
}

type JWTServiceOption func(*JWTService) error

func NewJWTService(options ...JWTServiceOption) (*JWTService, error) {
	service := JWTService{
		issuer:            "dynamic-webforms",
		signingMethod:     jwt.SigningMethodHS512,
		signingMethodAlg:  []string{jwt.SigningMethodHS512.Alg()},
		expiryTimeMinutes: 30,
		signingKey:        nil,
		parser:            nil,
	}

	for _, option := range options {
		if err := option(&service); err != nil {
			return nil, err
		}
	}

	if service.signingKey == nil {
		service.signingKey = make([]byte, 128)
		if _, err := rand.Read(service.signingKey); err != nil {
			return nil, fmt.Errorf("failed to generate signing key: %w", err)
		}
	}

	service.parser = jwt.NewParser(jwt.WithIssuer(service.issuer), jwt.WithExpirationRequired(), jwt.WithIssuedAt(),
		jwt.WithValidMethods(service.signingMethodAlg))

	return &service, nil
}

func WithIssuer(issuer string) JWTServiceOption {
	return func(s *JWTService) error {
		s.issuer = issuer
		return nil
	}
}

func WithExpiryTimeMinutes(expiryTimeMinutes int) JWTServiceOption {
	return func(s *JWTService) error {
		if expiryTimeMinutes <= 0 {
			return errors.New("expiryTimeMinutes must be greater than zero")
		}

		s.expiryTimeMinutes = expiryTimeMinutes
		return nil
	}
}

func WithSigningKey(signingKey []byte) JWTServiceOption {
	return func(s *JWTService) error {
		if signingKey == nil || len(signingKey) < 64 {
			return errors.New("signingKey must be greater than 64 bytes")
		}
		emptyKey := true
		for i := 0; i < len(signingKey); i++ {
			if signingKey[i] != 0 {
				emptyKey = false
				break
			}
		}

		if emptyKey {
			return errors.New("signingKey is empty")
		}

		s.signingKey = signingKey
		return nil
	}
}

func (j *JWTService) NewToken(userID uuid.UUID) (string, error) {
	currentTime := time.Now().UTC()

	randomID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to generate random ID: %w", err)
	}

	claims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    j.issuer,
			Subject:   userID.String(),
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(time.Minute * time.Duration(j.expiryTimeMinutes))),
			IssuedAt:  jwt.NewNumericDate(currentTime),
			ID:        randomID.String(),
		},
	}

	return jwt.NewWithClaims(j.signingMethod, claims).SignedString(j.signingKey)
}

func (j *JWTService) ValidateToken(tokenString string) (JWTClaims, error) {
	claims := &JWTClaims{}
	_, err := j.parser.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return j.signingKey, nil
	})

	if err != nil {
		return JWTClaims{}, err
	}

	return *claims, nil
}
