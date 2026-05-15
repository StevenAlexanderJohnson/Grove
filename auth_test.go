package grove_test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/StevenAlexanderJohnson/grove"
	"github.com/golang-jwt/jwt/v5"
)

type TestClaims struct {
	Email string `json:"email"`
	*jwt.RegisteredClaims
}

func testRSAKey(t *testing.T) *rsa.PrivateKey {
	t.Helper()

	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatalf("failed to generate RSA key: %v", err)
	}

	return key
}

func writePrivateKeyPEM(t *testing.T, key *rsa.PrivateKey) string {
	t.Helper()

	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("failed to marshal private key: %v", err)
	}

	path := filepath.Join(t.TempDir(), "private.pem")

	err = os.WriteFile(path, pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	}), 0600)
	if err != nil {
		t.Fatalf("failed to write private key: %v", err)
	}

	return path
}

func validConfig(t *testing.T, canEncrypt bool) grove.AuthenticatorConfig {
	t.Helper()

	var key *rsa.PrivateKey
	if canEncrypt {
		key = testRSAKey(t)
	}

	return grove.AuthenticatorConfig{
		CanEncrypt:    canEncrypt,
		JWEPrivateKey: key,
		Lifetime:      time.Hour,
		Issuer:        "Testing",
		Audience:      []string{"testing"},
		Key:           "secret",
	}
}

func validClaims() *TestClaims {
	return &TestClaims{
		Email: "testing@example.com",
		RegisteredClaims: &jwt.RegisteredClaims{
			Issuer:    "Testing",
			Audience:  []string{"testing"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func TestNewAuthenticatorConfig(t *testing.T) {
	key := testRSAKey(t)
	audience := []string{"testing"}

	got := grove.NewAuthenticatorConfig(
		true,
		key,
		time.Hour,
		"testing",
		audience,
		"key",
	)

	if !got.CanEncrypt {
		t.Fatalf("CanEncrypt = %t; want true", got.CanEncrypt)
	}
	if got.JWEPrivateKey != key {
		t.Fatalf("JWEPrivateKey does not match provided key")
	}
	if got.Lifetime != time.Hour {
		t.Fatalf("Lifetime = %s; want %s", got.Lifetime, time.Hour)
	}
	if got.Issuer != "testing" {
		t.Fatalf("Issuer = %s; want testing", got.Issuer)
	}
	if slices.Compare(got.Audience, audience) != 0 {
		t.Fatalf("Audience = %v; want %v", got.Audience, audience)
	}
	if got.Key != "key" {
		t.Fatalf("Key = %s; want key", got.Key)
	}
}

func TestAuthenticatorConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		config  grove.AuthenticatorConfig
		wantErr bool
	}{
		{
			name:    "valid with encryption",
			config:  validConfig(t, true),
			wantErr: false,
		},
		{
			name:    "valid without encryption",
			config:  validConfig(t, false),
			wantErr: false,
		},
		{
			name: "missing private key when encryption enabled",
			config: grove.AuthenticatorConfig{
				CanEncrypt: true,
				Lifetime:   time.Hour,
				Issuer:     "Testing",
				Audience:   []string{"testing"},
				Key:        "secret",
			},
			wantErr: true,
		},
		{
			name: "invalid lifetime",
			config: grove.AuthenticatorConfig{
				Lifetime: -time.Hour,
				Issuer:   "Testing",
				Audience: []string{"testing"},
				Key:      "secret",
			},
			wantErr: true,
		},
		{
			name: "missing issuer",
			config: grove.AuthenticatorConfig{
				Lifetime: time.Hour,
				Audience: []string{"testing"},
				Key:      "secret",
			},
			wantErr: true,
		},
		{
			name: "missing audience",
			config: grove.AuthenticatorConfig{
				Lifetime: time.Hour,
				Issuer:   "Testing",
				Key:      "secret",
			},
			wantErr: true,
		},
		{
			name: "missing key",
			config: grove.AuthenticatorConfig{
				Lifetime: time.Hour,
				Issuer:   "Testing",
				Audience: []string{"testing"},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr && err == nil {
				t.Fatalf("Validate() error = nil; want error")
			}
			if !tt.wantErr && err != nil {
				t.Fatalf("Validate() error = %v; want nil", err)
			}
		})
	}
}

func TestLoadAuthenticatorConfigFromEnvWithoutEncryption(t *testing.T) {
	t.Setenv("JWT_CAN_ENCRYPT", "false")
	t.Setenv("JWT_LIFETIME", "30")
	t.Setenv("JWT_ISSUER", "Testing")
	t.Setenv("JWT_AUDIENCE", "testing,admin")
	t.Setenv("JWT_SECRET", "secret")

	got, err := grove.LoadAuthenticatorConfigFromEnv()
	if err != nil {
		t.Fatalf("LoadAuthenticatorConfigFromEnv() error = %v; want nil", err)
	}

	if got.CanEncrypt {
		t.Fatalf("CanEncrypt = true; want false")
	}
	if got.JWEPrivateKey != nil {
		t.Fatalf("JWEPrivateKey = %v; want nil", got.JWEPrivateKey)
	}
	if got.Lifetime != 30*time.Minute {
		t.Fatalf("Lifetime = %s; want %s", got.Lifetime, 30*time.Minute)
	}
	if got.Issuer != "Testing" {
		t.Fatalf("Issuer = %s; want Testing", got.Issuer)
	}
	if slices.Compare(got.Audience, []string{"admin", "testing"}) != 0 {
		t.Fatalf("Audience = %v; want sorted [admin testing]", got.Audience)
	}
	if got.Key != "secret" {
		t.Fatalf("Key = %s; want secret", got.Key)
	}
}

func TestLoadAuthenticatorConfigFromEnvWithEncryption(t *testing.T) {
	key := testRSAKey(t)
	keyPath := writePrivateKeyPEM(t, key)

	t.Setenv("JWT_CAN_ENCRYPT", "true")
	t.Setenv("JWT_PRIVATE_KEY_PATH", keyPath)
	t.Setenv("JWT_LIFETIME", "30")
	t.Setenv("JWT_ISSUER", "Testing")
	t.Setenv("JWT_AUDIENCE", "testing")
	t.Setenv("JWT_SECRET", "secret")

	got, err := grove.LoadAuthenticatorConfigFromEnv()
	if err != nil {
		t.Fatalf("LoadAuthenticatorConfigFromEnv() error = %v; want nil", err)
	}

	if !got.CanEncrypt {
		t.Fatalf("CanEncrypt = false; want true")
	}
	if got.JWEPrivateKey == nil {
		t.Fatalf("JWEPrivateKey = nil; want RSA key")
	}
}

func TestNewAuthenticatorWithNilConfigShouldFail(t *testing.T) {
	got, err := grove.NewAuthenticator[*jwt.RegisteredClaims](nil)
	if err == nil {
		t.Fatalf("NewAuthenticator(nil) error = nil; want error")
	}
	if got != nil {
		t.Fatalf("NewAuthenticator(nil) authenticator = %v; want nil", got)
	}
}

func TestGenerateAndVerifyEncryptedToken(t *testing.T) {
	config := validConfig(t, true)

	auth, err := grove.NewAuthenticator[*TestClaims](&config)
	if err != nil {
		t.Fatalf("NewAuthenticator() error = %v; want nil", err)
	}

	claims := validClaims()
	token, err := auth.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v; want nil", err)
	}

	if parts := strings.Count(token, ".") + 1; parts != 5 {
		t.Fatalf("encrypted token has %d parts; want 5-part JWE", parts)
	}

	got, err := auth.VerifyToken(token, &TestClaims{})
	if err != nil {
		t.Fatalf("VerifyToken() error = %v; want nil", err)
	}

	if got.Issuer != config.Issuer {
		t.Fatalf("Issuer = %s; want %s", got.Issuer, config.Issuer)
	}
}

func TestGenerateAndVerifyUnencryptedToken(t *testing.T) {
	config := validConfig(t, false)

	auth, err := grove.NewAuthenticator[*TestClaims](&config)
	if err != nil {
		t.Fatalf("NewAuthenticator() error = %v; want nil", err)
	}

	claims := validClaims()
	token, err := auth.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v; want nil", err)
	}

	if parts := strings.Count(token, ".") + 1; parts != 3 {
		t.Fatalf("unencrypted token has %d parts; want 3-part JWT", parts)
	}

	got, err := auth.VerifyToken(token, &TestClaims{})
	if err != nil {
		t.Fatalf("VerifyToken() error = %v; want nil", err)
	}

	if got.Issuer != config.Issuer {
		t.Fatalf("Issuer = %s; want %s", got.Issuer, config.Issuer)
	}
}

func TestParseEncryptedToken(t *testing.T) {
	config := validConfig(t, true)

	auth, err := grove.NewAuthenticator[*TestClaims](&config)
	if err != nil {
		t.Fatalf("NewAuthenticator() error = %v; want nil", err)
	}

	token, err := auth.GenerateToken(validClaims())
	if err != nil {
		t.Fatalf("GenerateToken() error = %v; want nil", err)
	}

	parsed, err := auth.ParseToken(token, &TestClaims{
		RegisteredClaims: &jwt.RegisteredClaims{},
	})
	if err != nil {
		t.Fatalf("ParseToken() error = %v; want nil", err)
	}

	claims, ok := parsed.Claims.(*TestClaims)
	if !ok {
		t.Fatalf("Claims type = %T; want *TestClaims", parsed.Claims)
	}

	if claims.Email != "testing@example.com" {
		t.Fatalf("Email = %s; want testing@example.com", claims.Email)
	}
}

func TestParseUnencryptedToken(t *testing.T) {
	config := validConfig(t, false)

	auth, err := grove.NewAuthenticator[*TestClaims](&config)
	if err != nil {
		t.Fatalf("NewAuthenticator() error = %v; want nil", err)
	}

	token, err := auth.GenerateToken(validClaims())
	if err != nil {
		t.Fatalf("GenerateToken() error = %v; want nil", err)
	}

	parsed, err := auth.ParseToken(token, &TestClaims{
		RegisteredClaims: &jwt.RegisteredClaims{},
	})
	if err != nil {
		t.Fatalf("ParseToken() error = %v; want nil", err)
	}

	claims, ok := parsed.Claims.(*TestClaims)
	if !ok {
		t.Fatalf("Claims type = %T; want *TestClaims", parsed.Claims)
	}

	if claims.Email != "testing@example.com" {
		t.Fatalf("Email = %s; want testing@example.com", claims.Email)
	}
}

func TestParseTokenWithInvalidTokenShouldFail(t *testing.T) {
	config := validConfig(t, false)

	auth, err := grove.NewAuthenticator[*TestClaims](&config)
	if err != nil {
		t.Fatalf("NewAuthenticator() error = %v; want nil", err)
	}

	_, err = auth.ParseToken("not-a-valid-token", &TestClaims{
		RegisteredClaims: &jwt.RegisteredClaims{},
	})
	if err == nil {
		t.Fatalf("ParseToken() error = nil; want error")
	}
}

func TestParseEncryptedTokenWithInvalidJWEShouldFail(t *testing.T) {
	config := validConfig(t, true)

	auth, err := grove.NewAuthenticator[*TestClaims](&config)
	if err != nil {
		t.Fatalf("NewAuthenticator() error = %v; want nil", err)
	}

	_, err = auth.ParseToken("not-a-valid-jwe", &TestClaims{
		RegisteredClaims: &jwt.RegisteredClaims{},
	})
	if err == nil {
		t.Fatalf("ParseToken() error = nil; want error")
	}
}

func TestVerifyTokenWithWrongIssuerShouldFail(t *testing.T) {
	config := validConfig(t, false)

	auth, err := grove.NewAuthenticator[*TestClaims](&config)
	if err != nil {
		t.Fatalf("NewAuthenticator() error = %v; want nil", err)
	}

	claims := validClaims()
	claims.Issuer = "WrongIssuer"

	token, err := auth.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v; want nil", err)
	}

	_, err = auth.VerifyToken(token, &TestClaims{
		RegisteredClaims: &jwt.RegisteredClaims{},
	})
	if err == nil {
		t.Fatalf("VerifyToken() error = nil; want issuer error")
	}
}

func TestVerifyTokenWithWrongAudienceShouldFail(t *testing.T) {
	config := validConfig(t, false)

	auth, err := grove.NewAuthenticator[*TestClaims](&config)
	if err != nil {
		t.Fatalf("NewAuthenticator() error = %v; want nil", err)
	}

	claims := validClaims()
	claims.Audience = []string{"wrong-audience"}

	token, err := auth.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v; want nil", err)
	}

	_, err = auth.VerifyToken(token, &TestClaims{
		RegisteredClaims: &jwt.RegisteredClaims{},
	})
	if err == nil {
		t.Fatalf("VerifyToken() error = nil; want audience error")
	}
}

func TestVerifyTokenWithExpiredTokenShouldFail(t *testing.T) {
	config := validConfig(t, false)

	auth, err := grove.NewAuthenticator[*TestClaims](&config)
	if err != nil {
		t.Fatalf("NewAuthenticator() error = %v; want nil", err)
	}

	claims := validClaims()
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Hour))

	token, err := auth.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v; want nil", err)
	}

	_, err = auth.VerifyToken(token, &TestClaims{
		RegisteredClaims: &jwt.RegisteredClaims{},
	})
	if err == nil {
		t.Fatalf("VerifyToken() error = nil; want expiration error")
	}
}

func TestVerifyTokenWithoutExpirationShouldFail(t *testing.T) {
	config := validConfig(t, false)

	auth, err := grove.NewAuthenticator[*TestClaims](&config)
	if err != nil {
		t.Fatalf("NewAuthenticator() error = %v; want nil", err)
	}

	claims := validClaims()
	claims.ExpiresAt = nil

	token, err := auth.GenerateToken(claims)
	if err != nil {
		t.Fatalf("GenerateToken() error = %v; want nil", err)
	}

	_, err = auth.VerifyToken(token, &TestClaims{
		RegisteredClaims: &jwt.RegisteredClaims{},
	})
	if err == nil {
		t.Fatalf("VerifyToken() error = nil; want missing expiration error")
	}
}

func TestLoadAuthenticatorConfigFromEnvWithInvalidCanEncryptShouldFail(t *testing.T) {
	t.Setenv("JWT_CAN_ENCRYPT", "not-a-bool")
	t.Setenv("JWT_ISSUER", "Testing")
	t.Setenv("JWT_AUDIENCE", "testing")
	t.Setenv("JWT_SECRET", "secret")

	_, err := grove.LoadAuthenticatorConfigFromEnv()
	if err == nil {
		t.Fatalf("LoadAuthenticatorConfigFromEnv() error = nil; want error")
	}
}

func TestLoadAuthenticatorConfigFromEnvWithInvalidLifetimeShouldFail(t *testing.T) {
	t.Setenv("JWT_CAN_ENCRYPT", "false")
	t.Setenv("JWT_LIFETIME", "not-an-int")
	t.Setenv("JWT_ISSUER", "Testing")
	t.Setenv("JWT_AUDIENCE", "testing")
	t.Setenv("JWT_SECRET", "secret")

	_, err := grove.LoadAuthenticatorConfigFromEnv()
	if err == nil {
		t.Fatalf("LoadAuthenticatorConfigFromEnv() error = nil; want error")
	}
}

func TestLoadAuthenticatorConfigFromEnvMissingIssuerShouldFail(t *testing.T) {
	t.Setenv("JWT_CAN_ENCRYPT", "false")
	t.Setenv("JWT_AUDIENCE", "testing")
	t.Setenv("JWT_SECRET", "secret")

	_, err := grove.LoadAuthenticatorConfigFromEnv()
	if err == nil {
		t.Fatalf("LoadAuthenticatorConfigFromEnv() error = nil; want error")
	}
}

func TestLoadAuthenticatorConfigFromEnvMissingAudienceShouldFail(t *testing.T) {
	t.Setenv("JWT_CAN_ENCRYPT", "false")
	t.Setenv("JWT_ISSUER", "Testing")
	t.Setenv("JWT_SECRET", "secret")

	_, err := grove.LoadAuthenticatorConfigFromEnv()
	if err == nil {
		t.Fatalf("LoadAuthenticatorConfigFromEnv() error = nil; want error")
	}
}

func TestLoadAuthenticatorConfigFromEnvMissingSecretShouldFail(t *testing.T) {
	t.Setenv("JWT_CAN_ENCRYPT", "false")
	t.Setenv("JWT_ISSUER", "Testing")
	t.Setenv("JWT_AUDIENCE", "testing")

	_, err := grove.LoadAuthenticatorConfigFromEnv()
	if err == nil {
		t.Fatalf("LoadAuthenticatorConfigFromEnv() error = nil; want error")
	}
}

func TestLoadAuthenticatorConfigFromEnvMissingPrivateKeyPathShouldFail(t *testing.T) {
	t.Setenv("JWT_CAN_ENCRYPT", "true")
	t.Setenv("JWT_ISSUER", "Testing")
	t.Setenv("JWT_AUDIENCE", "testing")
	t.Setenv("JWT_SECRET", "secret")

	_, err := grove.LoadAuthenticatorConfigFromEnv()
	if err == nil {
		t.Fatalf("LoadAuthenticatorConfigFromEnv() error = nil; want error")
	}
}
