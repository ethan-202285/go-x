package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

const gcmStandardNonceSize = 12

// NewGCM 新建
func NewGCM(key []byte) *GCM {
	return &GCM{secretKey: key}
}

// GCM GCM加密
type GCM struct {
	secretKey []byte
}

// Encrypt 加密
func (a *GCM) Encrypt(plaintext []byte) []byte {
	block, err := aes.NewCipher(a.secretKey)
	if err != nil {
		panic(err)
	}

	// Never use more than 2^32 random nonces with a given key
	// because of the risk of a repeat.
	nonce := make([]byte, gcmStandardNonceSize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	// This encrypts `content` and appends the result to the nonce.
	ciphertext := aesgcm.Seal(nonce[:], nonce, plaintext, nil)
	return ciphertext
}

// Decrypt 解密
func (a *GCM) Decrypt(encrypted []byte) (plaintext []byte, err error) {
	if len(encrypted) <= gcmStandardNonceSize {
		return nil, fmt.Errorf("error ciphertext size (<=%d)", gcmStandardNonceSize)
	}

	var nonce []byte
	nonce = encrypted[:gcmStandardNonceSize]

	block, err := aes.NewCipher(a.secretKey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err = aesgcm.Open(nil, nonce, encrypted[gcmStandardNonceSize:], nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}
