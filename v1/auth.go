package grove

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/golang-jwt/jwt/v5"
)

type AuthenticatorConfig struct {
	JWEPrivateKey *rsa.PrivateKey
	Lifetime      time.Duration
	Issuer        string
	Audience      []string
	Key           string
}

func NewAuthenticatorConfig(jwePrivateKey *rsa.PrivateKey, lifetime time.Duration, issuer string, audience []string, key string) *AuthenticatorConfig {
	return &AuthenticatorConfig{
		JWEPrivateKey: jwePrivateKey,
		Lifetime:      lifetime,
		Issuer:        issuer,
		Audience:      audience,
		Key:           key,
	}
}

func (config *AuthenticatorConfig) Validate() error {
	if config.JWEPrivateKey == nil {
		return fmt.Errorf("JWEPrivateKey is required")
	}
	if config.Lifetime <= 0 {
		return fmt.Errorf("lifetime must be greater than zero")
	}
	if config.Issuer == "" {
		return fmt.Errorf("issuer is required")
	}
	if len(config.Audience) == 0 {
		return fmt.Errorf("audience must contain at least one value")
	}
	if config.Key == "" {
		return fmt.Errorf("key is required")
	}
	return nil
}

func LoadAuthenticatorConfigFromEnv() (*AuthenticatorConfig, error) {
	jwePem, err := os.ReadFile(os.Getenv("JWT_PRIVATE_KEY_PATH"))
	if err != nil {
		return nil, fmt.Errorf("an error occurred while reading JWE private key: %v", err)
	}
	block, _ := pem.Decode(jwePem)
	if block == nil {
		return nil, fmt.Errorf("an error occurred while decoding JWE private key")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while parsing JWE private key: %v", err)
	}
	rsaKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("an error occurred while converting JWE private key to RSA: %v", err)
	}

	lifetimeSetting := os.Getenv("JWT_LIFETIME")
	var lifetime time.Duration
	if lifetimeSetting != "" {
		l, err := strconv.ParseInt(lifetimeSetting, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("an error occurred while loading jwt lifetime: %v", err)
		}
		lifetime = time.Duration(time.Duration(l) * time.Minute)
	}
	if lifetime == 0 {
		lifetime = time.Duration(2 * time.Minute)
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		return nil, fmt.Errorf("JWT_ISSUER was not set")
	}

	audienceSetting := os.Getenv("JWT_AUDIENCE")
	if audienceSetting == "" {
		return nil, fmt.Errorf("JWT_AUDIENCE was not set")
	}
	audience := strings.Split(audienceSetting, ",")
	slices.Sort(audience)

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET was not set")
	}

	jwtConfig := &AuthenticatorConfig{
		Lifetime:      time.Duration(2 * time.Minute),
		Issuer:        issuer,
		Audience:      audience,
		Key:           jwtSecret,
		JWEPrivateKey: rsaKey,
	}

	return jwtConfig, nil

}

func NewAuthenticator[T jwt.Claims](config *AuthenticatorConfig) *Authenticator[T] {
	return &Authenticator[T]{
		AuthenticatorConfig: config,
	}
}

type Authenticator[T jwt.Claims] struct {
	*AuthenticatorConfig
}

func (a *Authenticator[T]) encryptToken(token string) (string, error) {
	encryptor, err := jose.NewEncrypter(jose.A128GCM, jose.Recipient{Algorithm: jose.RSA_OAEP, Key: &a.JWEPrivateKey.PublicKey}, nil)
	if err != nil {
		return "", fmt.Errorf("an error occurred while creating encrypter: %v", err)
	}

	jweObject, err := encryptor.Encrypt([]byte(token))
	if err != nil {
		return "", fmt.Errorf("an error occurred while encrypting JWT: %v", err)
	}

	compact, err := jweObject.CompactSerialize()
	if err != nil {
		return "", fmt.Errorf("an error occurred while serializing JWE: %v", err)
	}

	return compact, nil
}

func (a *Authenticator[T]) decryptToken(encrypted string) (string, error) {
	parsedCompact, err := jose.ParseEncrypted(encrypted, []jose.KeyAlgorithm{jose.RSA_OAEP}, []jose.ContentEncryption{jose.A128GCM})
	if err != nil {
		return "", fmt.Errorf("an error occurred while parsing JWE: %v", err)
	}

	tokenBytes, err := parsedCompact.Decrypt(a.JWEPrivateKey)
	if err != nil {
		return "", fmt.Errorf("an error occurred while decrypting JWE: %v", err)
	}

	return string(tokenBytes), nil
}

// GenerateToken creates a new JWT token with the provided claims, signs it, and encrypts it.
// The token is signed using the configured key and encrypted using JWE with RSA-OAEP
// and A128GCM.
// The generated token is suitable for use in authentication and authorization processes.
func (a *Authenticator[T]) GenerateToken(claims T) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedString, err := token.SignedString([]byte(a.Key))
	if err != nil {
		return "", fmt.Errorf("an error occurred while signing jwt: %v", err)
	}
	encrypted, err := a.encryptToken(signedString)
	if err != nil {
		return "", fmt.Errorf("an error occurred while encrypting jwt: %v", err)
	}

	return encrypted, nil
}

// ParseToken decrypts the token and parses it into the provided claims.
// It does not validate the claims, allowing for custom validation logic to be applied later.
func (a *Authenticator[T]) ParseToken(token string, claims T) (*jwt.Token, error) {
	decryptedToken, err := a.decryptToken(token)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while decrypting token: %v", err)
	}
	// We do not create a list of options because we are disabling all validation then doing manual validation.
	parsedToken, err := jwt.ParseWithClaims(
		decryptedToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(a.Key), nil
		},
		jwt.WithoutClaimsValidation(),
	)
	if err != nil {
		return nil, err
	}

	return parsedToken, nil
}

// VerifyToken decrypts the token, parses it, and validates the claims.
// It checks the audience and issuer against the configured values.
// If the token is valid, it returns the claims; otherwise, it returns an error.
// This method is used to ensure that the token is valid and can be trusted for authentication.
func (a *Authenticator[T]) VerifyToken(token string, claims T) (T, error) {
	decryptedToken, err := a.decryptToken(token)
	if err != nil {
		return claims, fmt.Errorf("an error occurred while decrypting token while verifying: %v", err)
	}

	parserOptions := make([]jwt.ParserOption, 0)
	for i := range a.Audience {
		parserOptions = append(parserOptions, jwt.WithAudience(a.Audience[i]))
	}
	parserOptions = append(parserOptions, jwt.WithIssuer(a.Issuer))
	parserOptions = append(parserOptions, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	parserOptions = append(parserOptions, jwt.WithExpirationRequired())

	parsedToken, err := jwt.ParseWithClaims(
		decryptedToken,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(a.Key), nil
		},
		parserOptions...,
	)

	if err != nil {
		return claims, fmt.Errorf("an error occurred while parsing JWT: %v", err)
	}
	if !parsedToken.Valid {
		return claims, fmt.Errorf("token is not valid")
	}
	return parsedToken.Claims.(T), nil
}
