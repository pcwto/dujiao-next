package config

import (
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func TestLoadCORSAllowedOriginsFromCommaSeparatedEnv(t *testing.T) {
	viper.Reset()
	t.Cleanup(viper.Reset)

	t.Setenv("CORS_ALLOWED_ORIGINS", "https://shop.example.com,https://admin.example.com")

	cfg := Load()
	want := []string{"https://shop.example.com", "https://admin.example.com"}
	if !reflect.DeepEqual(cfg.CORS.AllowedOrigins, want) {
		t.Fatalf("allowed origins mismatch, want %#v got %#v", want, cfg.CORS.AllowedOrigins)
	}
}
