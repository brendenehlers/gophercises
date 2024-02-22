package encrypt

import (
	"bytes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	key := []byte(`my super secret key`)
	plaintext := []byte("Hello, World!")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Errorf("Error encrypting: %v", err)
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Errorf("Error decrypting: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted text doesn't match original plaintext")
	}
}

func TestEncryptDecryptEmpty(t *testing.T) {
	key := []byte(`my super secret key`)
	plaintext := []byte("")

	ciphertext, err := Encrypt(key, plaintext)
	if err != nil {
		t.Errorf("Error encrypting: %v", err)
	}

	decrypted, err := Decrypt(key, ciphertext)
	if err != nil {
		t.Errorf("Error decrypting: %v", err)
	}

	if !bytes.Equal(plaintext, decrypted) {
		t.Errorf("Decrypted text doesn't match original plaintext")
	}
}
