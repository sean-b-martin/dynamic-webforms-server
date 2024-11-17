package auth

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewJWTService(t *testing.T) {
	service, err := NewJWTService()
	assert.NotNil(t, service)
	assert.NoError(t, err)
	assert.NotEmpty(t, service.signingKey)

	service2, err := NewJWTService()
	assert.NotNil(t, service2)
	assert.NoError(t, err)
	assert.NotEmpty(t, service2.signingKey)
	assert.NotEqual(t, service.signingKey, service2.signingKey)
}

func TestJWTService_NewToken(t *testing.T) {
	service, _ := NewJWTService()
	userID, _ := uuid.NewUUID()
	token, err := service.NewToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	token2, _ := service.NewToken(userID)
	assert.NotEmpty(t, token2)
	assert.NotEqual(t, token, token2)
}

func TestWithExpiryTimeMinutes(t *testing.T) {
	type args struct {
		expiryTimeMinutes int
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
	}{
		{name: "negative minutes", args: args{expiryTimeMinutes: -5}, wantError: true},
		{name: "zero minutes", args: args{expiryTimeMinutes: 0}, wantError: true},
		{name: "valid time", args: args{expiryTimeMinutes: 1}, wantError: false},
		{name: "valid time", args: args{expiryTimeMinutes: 1000}, wantError: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewJWTService(WithExpiryTimeMinutes(tt.args.expiryTimeMinutes))
			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, service)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.expiryTimeMinutes, service.expiryTimeMinutes)
			}
		})
	}
}

func TestWithIssuer(t *testing.T) {
	type args struct {
		issuer string
	}
	tests := []struct {
		name string
		args args
	}{
		{name: "empty name", args: args{issuer: ""}},
		{name: "not empty name", args: args{issuer: "test1"}},
		{name: "not empty name2", args: args{issuer: "test2"}},
		{name: "long name", args: args{issuer: "this is a long issuer name"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewJWTService(WithIssuer(tt.args.issuer))
			assert.NoError(t, err)
			assert.Equal(t, tt.args.issuer, service.issuer)
		})
	}
}

func TestWithSigningKey(t *testing.T) {
	type args struct {
		signingKey []byte
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
	}{
		{name: "no key", args: args{[]byte("")}, wantError: true},
		{name: "empty key", args: args{make([]byte, 64)}, wantError: true},
		{name: "valid key", args: args{[]byte("test-valid-and-secure-long-signing-key-greater-than-64-bytes!!!!")}, wantError: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewJWTService(WithSigningKey(tt.args.signingKey))
			if tt.wantError {
				assert.Error(t, err)
				assert.Empty(t, service)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.args.signingKey, service.signingKey)
			}
		})
	}
}

func TestJWTService_ValidateTokenInvalid(t *testing.T) {
	service, err := NewJWTService(WithSigningKey([]byte("test-valid-and-secure-long-signing-key-greater-than-64-bytes!!!!")))
	assert.NoError(t, err)
	assert.NotNil(t, service)

	token, err := service.NewToken(uuid.New())
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// valid token
	claims, err := service.ValidateToken(token)
	assert.NoError(t, err)
	assert.NotEmpty(t, claims)
	fmt.Println(token)

	// invalid tokens

	// expired
	claims, err = service.ValidateToken("eyJhbGciOiJIUzUxMiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJkeW5hbWljLXdlYmZvcm1zIiwic3ViIjoiMWU4MmVlY2EtNjkyOC00OWNkLWEwZTMtYzAxNjA1NDFjMjI3IiwiZXhwIjoxNzMxNzk5NDc5LCJpYXQiOjE3MzE3OTc2NzksImp0aSI6ImIzNTE3YTM4LTI0NDEtNGI2My04YTA4LTAzZGIyMmNiYjYyZSJ9.46gWFTwcufwyCrE1XkOkRm769BuxJr92EKPiMbVWCooKwdhDyLRoKtpRUjVa7KQtRAJKxux6RhD7OQwtLO3itA")
	assert.ErrorIs(t, err, jwt.ErrTokenExpired)
	assert.Empty(t, claims)

	// invalid signing key
	service2, _ := NewJWTService(WithSigningKey([]byte("test-valid-and-secure-long-signing-key-greater-than-64-bytes-but-different!!!!")))
	token, err = service2.NewToken(uuid.New())
	claims, err = service.ValidateToken(token)
	assert.ErrorIs(t, err, jwt.ErrTokenSignatureInvalid)
	assert.Empty(t, claims)
}
