package config

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	testConfigFetch, err := Init("./fixture", "test")

	if err != nil {
		t.Errorf("Init() error = %v, wantErr %v", err, testConfigFetch)

		return
	}

	testConfigWant := &Config{
		Environment: "test",
		HTTP: HTTPConfig{
			Host:              "local.test",
			Port:              "4051",
			ReadTimeout:       20 * time.Second,
			WriteTimeOut:      20 * time.Second,
			MaxHeaderMegabyte: 10,
		},
		Auth: AuthConfig{
			PasswordSalt:           "salt-test",
			VerificationCodeLength: 150,
			JWT: JWTConfig{
				SecretKey: "secret-key-test",
				AccessToken: TokenConfig{
					PrivateKey: "test-2",
					PublicKey:  "test-3",
					MaxAge:     15 * time.Minute,
				},
				RefreshToken: TokenConfig{
					PrivateKey: "test-2",
					PublicKey:  "test-3",
					MaxAge:     60 * time.Minute,
				},
			},
		},
	}

	fmt.Printf("testConfigFetch: %#v\n", testConfigFetch)
	fmt.Printf("testConfigWant:  %#v\n", testConfigWant)

	if !reflect.DeepEqual(testConfigFetch, testConfigWant) {
		t.Errorf("Init() got = %v, want %v", testConfigFetch, testConfigWant)
	}

}
