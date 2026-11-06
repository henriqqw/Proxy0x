package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func (a *App) keyFile() string {
	return filepath.Join(a.dataDir(), ".key")
}

func (a *App) getOrCreateKey() ([]byte, error) {
	os.MkdirAll(a.dataDir(), 0o700)

	path := a.keyFile()
	key, err := os.ReadFile(path)
	if err == nil && len(key) == 32 {
		return key, nil
	}

	key = make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return nil, err
	}
	if err := os.WriteFile(path, key, 0o600); err != nil {
		return nil, err
	}
	return key, nil
}

func (a *App) encrypt(plaintext string) (string, error) {
	if plaintext == "" {
		return "", nil
	}
	key, err := a.getOrCreateKey()
	if err != nil {
		return plaintext, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return plaintext, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return plaintext, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return plaintext, err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return "enc:" + base64.StdEncoding.EncodeToString(ciphertext), nil
}

func (a *App) decrypt(encoded string) (string, error) {
	if encoded == "" {
		return "", nil
	}
	if len(encoded) < 5 || encoded[:4] != "enc:" {
		return encoded, nil
	}

	key, err := a.getOrCreateKey()
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(encoded[4:])
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

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
