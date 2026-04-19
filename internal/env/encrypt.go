package env

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
)

type EncryptResult struct {
	Encrypted map[string]string
	Skipped   []string
	Total     int
}

func (r EncryptResult) Format() string {
	out := fmt.Sprintf("Encrypted: %d, Skipped: %d\n", len(r.Encrypted), len(r.Skipped))
	for _, k := range r.Skipped {
		out += fmt.Sprintf("  skipped: %s\n", k)
	}
	return out
}

// Encrypt encrypts secret values in entries using AES-GCM with the given key.
// Non-secret keys are skipped unless keys is non-empty (then only those keys are encrypted).
func Encrypt(entries []Entry, passKey string, keys []string) (EncryptResult, error) {
	keySet := make(map[string]bool)
	for _, k := range keys {
		keySet[k] = true
	}

	aesKey, err := deriveKey(passKey)
	if err != nil {
		return EncryptResult{}, err
	}

	result := EncryptResult{Encrypted: make(map[string]string), Total: len(entries)}

	for _, e := range entries {
		shouldEncrypt := (len(keys) > 0 && keySet[e.Key]) || (len(keys) == 0 && isSecret(e.Key))
		if !shouldEncrypt {
			result.Skipped = append(result.Skipped, e.Key)
			continue
		}
		enc, err := aesEncrypt(aesKey, e.Value)
		if err != nil {
			return EncryptResult{}, fmt.Errorf("encrypt %s: %w", e.Key, err)
		}
		result.Encrypted[e.Key] = enc
	}
	return result, nil
}

func deriveKey(pass string) ([]byte, error) {
	if len(pass) == 0 {
		return nil, errors.New("encryption key must not be empty")
	}
	key := make([]byte, 32)
	copy(key, []byte(pass))
	return key, nil
}

func aesEncrypt(key []byte, plaintext string) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	sealed := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return hex.EncodeToString(sealed), nil
}

func AesDecrypt(key []byte, cipherHex string) (string, error) {
	data, err := hex.DecodeString(cipherHex)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ns := gcm.NonceSize()
	if len(data) < ns {
		return "", errors.New("ciphertext too short")
	}
	plain, err := gcm.Open(nil, data[:ns], data[ns:], nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
