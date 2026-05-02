package main

import (
	"testing"

	"github.com/dujiao-next/internal/config"
)

func TestValidateProductionSecretsRejectsWeakOrRepeatedSecrets(t *testing.T) {
	strongAdmin := "admin-secret-32-bytes-minimum-value"
	strongUser := "user-secret-32-bytes-minimum-value!"
	strongApp := "app-secret-32-bytes-minimum-value!!"

	tests := []struct {
		name string
		cfg  config.Config
		want bool
	}{
		{
			name: "all secrets strong and unique",
			cfg: config.Config{
				JWT:     config.JWTConfig{SecretKey: strongAdmin},
				UserJWT: config.JWTConfig{SecretKey: strongUser},
				App:     config.AppConfig{SecretKey: strongApp},
			},
			want: true,
		},
		{
			name: "default user jwt secret",
			cfg: config.Config{
				JWT:     config.JWTConfig{SecretKey: strongAdmin},
				UserJWT: config.JWTConfig{SecretKey: "user-change-me-in-production"},
				App:     config.AppConfig{SecretKey: strongApp},
			},
			want: false,
		},
		{
			name: "default app secret",
			cfg: config.Config{
				JWT:     config.JWTConfig{SecretKey: strongAdmin},
				UserJWT: config.JWTConfig{SecretKey: strongUser},
				App:     config.AppConfig{SecretKey: "change-me-32-byte-secret-key!!"},
			},
			want: false,
		},
		{
			name: "repeated jwt secrets",
			cfg: config.Config{
				JWT:     config.JWTConfig{SecretKey: strongAdmin},
				UserJWT: config.JWTConfig{SecretKey: strongAdmin},
				App:     config.AppConfig{SecretKey: strongApp},
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateProductionSecrets(&tt.cfg)
			if tt.want && err != nil {
				t.Fatalf("validateProductionSecrets() returned error: %v", err)
			}
			if !tt.want && err == nil {
				t.Fatal("validateProductionSecrets() expected error")
			}
		})
	}
}
