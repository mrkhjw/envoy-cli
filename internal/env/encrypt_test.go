package env

import (
	"testing"
)

func TestEncrypt_SecretKeys(t *testing.T) {
	entries := []Entry{
		{Key: "API_SECRET", Value: "mysecret"},
		{Key: "APP_NAME", Value: "myapp"},
	}
	res, err := Encrypt(entries, "passphrase1234567890123456789012", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Encrypted["API_SECRET"]; !ok {
		t.Error("expected API_SECRET to be encrypted")
	}
	if _, ok := res.Encrypted["APP_NAME"]; ok {
		t.Error("expected APP_NAME to be skipped")
	}
}

func TestEncrypt_SpecificKeys(t *testing.T) {
	entries := []Entry{
		{Key: "APP_NAME", Value: "myapp"},
		{Key: "DB_HOST", Value: "localhost"},
	}
	res, err := Encrypt(entries, "passphrase1234567890123456789012", []string{"APP_NAME"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := res.Encrypted["APP_NAME"]; !ok {
		t.Error("expected APP_NAME to be encrypted")
	}
	if _, ok := res.Encrypted["DB_HOST"]; ok {
		t.Error("expected DB_HOST to be skipped")
	}
}

func TestEncrypt_EmptyPassKey(t *testing.T) {
	entries := []Entry{{Key: "API_SECRET", Value: "val"}}
	_, err := Encrypt(entries, "", nil)
	if err == nil {
		t.Error("expected error for empty passkey")
	}
}

func TestEncrypt_DecryptRoundtrip(t *testing.T) {
	entries := []Entry{{Key: "API_TOKEN", Value: "super-secret-value"}}
	pass := "passphrase1234567890123456789012"
	res, err := Encrypt(entries, pass, nil)
	if err != nil {
		t.Fatalf("encrypt error: %v", err)
	}
	enc := res.Encrypted["API_TOKEN"]
	key, _ := deriveKey(pass)
	plain, err := AesDecrypt(key, enc)
	if err != nil {
		t.Fatalf("decrypt error: %v", err)
	}
	if plain != "super-secret-value" {
		t.Errorf("expected 'super-secret-value', got %q", plain)
	}
}

func TestEncryptResult_Format(t *testing.T) {
	res := EncryptResult{
		Encrypted: map[string]string{"API_SECRET": "abc"},
		Skipped:   []string{"APP_NAME"},
		Total:     2,
	}
	out := res.Format()
	if out == "" {
		t.Error("expected non-empty format output")
	}
}
