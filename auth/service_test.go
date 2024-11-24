package auth

import (
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestNewPasswordService(t *testing.T) {
	type args struct {
		cost int
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
	}{
		{name: "negative cost", args: args{cost: -1}, wantError: true},
		{name: "zero cost", args: args{cost: 0}, wantError: true},
		{name: "wantError cost", args: args{cost: bcrypt.DefaultCost}, wantError: false},
		{name: "too large cost", args: args{cost: bcrypt.MaxCost + 1}, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewPasswordService(tt.args.cost)
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, service)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, service)
				assert.Equal(t, tt.args.cost, service.cost)
			}
		})
	}
}

func TestPasswordService_Hash(t *testing.T) {
	type fields struct {
		cost int
	}
	type args struct {
		password string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{name: "empty password", fields: fields{cost: bcrypt.DefaultCost}, args: args{password: ""}},
		{name: "normal password", fields: fields{cost: bcrypt.DefaultCost}, args: args{password: "abc"}},
		{name: "long password", fields: fields{cost: bcrypt.DefaultCost}, args: args{password: "A#4zT!fW9Pq@&eR6m^X*o(3K=+Lh|~78d%_CnG}"}},
		{name: "password with unicode", fields: fields{cost: bcrypt.DefaultCost}, args: args{password: "🔒🔒unicode-password🔒🔒"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewPasswordService(tt.fields.cost)
			assert.NoError(t, err)
			got, err := service.HashPassword(tt.args.password)
			assert.NoError(t, err)
			assert.NotEmpty(t, got)
			assert.Len(t, got, 60)
		})
	}
}

func TestPasswordService_VerifyPassword(t *testing.T) {
	type fields struct {
		cost int
	}
	type args struct {
		hash     string
		password string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantError bool
	}{
		{name: "password", fields: fields{cost: bcrypt.DefaultCost}, args: args{hash: "$2a$10$nIhWUx57Dc44tOxXUN.dB.ihWQY4wTVuSTrIDGpFAgnqgEGj79qkK", password: "A#4zT!fW9Pq@&eR6m^X*o(3K=+Lh|~78d%_CnG}"}, wantError: false},
		{name: "unicode password", fields: fields{cost: bcrypt.DefaultCost}, args: args{hash: "$2a$10$D9S359yaos1yLuUfbJH5iemxxOmgPX/5YrfA3.R72ppkIQZFeLmjq", password: "🔒🔒unicode-password🔒🔒"}, wantError: false},
		{name: "high cost", fields: fields{cost: 15}, args: args{hash: "$2a$15$lsWyUkdmgXLO7UkWIt1dauxF8H7O4VOo2YVyc4uKLJhxrumLlOyVO", password: "another!password1234"}, wantError: false},
		{name: "password, invalid", fields: fields{cost: bcrypt.DefaultCost}, args: args{hash: "$2a$10$nIhWUx57Dc44tOxXUN.dB.ihWQY4wTVuSTrIDGpFAgnqgEGj79qkK", password: "A#4zT!fW9Pq@&eR6m^X*o(3K=+Lh|~78d%_Cnfx}"}, wantError: true},
		{name: "unicode password, invalid", fields: fields{cost: bcrypt.DefaultCost}, args: args{hash: "$2a$10$D9S359yaos1yLuUfbJH5iemxxOmgPX/5YrfA3.R72ppkIQZFeLmjq", password: "🔒🔒unicode-password🔒!"}, wantError: true},
		{name: "high cost, invalid", fields: fields{cost: 15}, args: args{hash: "$2a$15$lsWyUkdmgXLO7UkWIt1dauxF8H7O4VOo2YVyc4uKLJhxrumLlOyVO", password: "another!pas"}, wantError: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service, err := NewPasswordService(tt.fields.cost)
			assert.NoError(t, err)

			if tt.wantError {
				assert.Error(t, service.VerifyPassword(tt.args.hash, tt.args.password))
			} else {
				assert.NoError(t, service.VerifyPassword(tt.args.hash, tt.args.password))
			}
		})
	}
}
