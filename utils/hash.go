package utils

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

type Argon2Params struct {
	Memory      uint32
	Iterations  uint32
	Parallelism uint8
	SaltLength  uint32
	KeyLength   uint32
}

var DefaultParams = &Argon2Params{
	Memory:      64 * 1024, // 64 MB
	Iterations:  3,
	Parallelism: 2,
	SaltLength:  16,
	KeyLength:   32,
}

func GenerateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	return b, err
}

// HashPassword returns an encoded hash that includes parameters + salt + key.
func HashPassword(password string, p *Argon2Params) (string, error) {
	if p == nil {
		p = DefaultParams
	}

	salt, err := GenerateRandomBytes(p.SaltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.Iterations, p.Memory, p.Parallelism, p.KeyLength)

	// Encode to a string like: argon2id$mem=65536,iter=3,par=2$salt$hash
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encoded := fmt.Sprintf("argon2id$mem=%d,iter=%d,par=%d$%s$%s",
		p.Memory, p.Iterations, p.Parallelism, b64Salt, b64Hash)

	return encoded, nil
}

func HashPasswordDefault(password string) (string, error) {
	return HashPassword(password, DefaultParams)
}

// ComparePasswordAndHash verifies password against encoded hash.
func ComparePasswordAndHash(password, encodedHash string) (bool, error) {
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 4 {
		return false, errors.New("invalid encoded hash format")
	}

	if parts[0] != "argon2id" {
		return false, fmt.Errorf("unsupported algorithm: %s", parts[0])
	}

	var mem uint32
	var iter uint32
	var par uint8
	_, err := fmt.Sscanf(parts[1], "mem=%d,iter=%d,par=%d", &mem, &iter, &par)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[2])
	if err != nil {
		return false, err
	}

	hash, err := base64.RawStdEncoding.DecodeString(parts[3])
	if err != nil {
		return false, err
	}

	keyLen := uint32(len(hash))

	computed := argon2.IDKey([]byte(password), salt, iter, mem, par, keyLen)

	// Constant-time compare
	if subtle.ConstantTimeCompare(hash, computed) == 1 {
		return true, nil
	}
	return false, nil
}
