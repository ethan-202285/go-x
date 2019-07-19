package crypto

import (
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

// NewNaCL 新建
func NewNaCL(key []byte) *NaCL {
	// Load your secret key from a safe place and reuse it across multiple
	// Seal calls. If you want to convert a passphrase to a key, use a suitable
	// package like bcrypt or scrypt.
	var secretKey [32]byte
	copy(secretKey[:], key)
	return &NaCL{secretKey: secretKey}
}

// NaCL GCM加密
type NaCL struct {
	secretKey [32]byte
}

// Encrypt 加密
func (s *NaCL) Encrypt(plaintext []byte) []byte {
	// You must use a different nonce for each message you encrypt with the
	// same key. Since the nonce here is 192 bits long, a random value
	// provides a sufficiently small probability of repeats.
	// nacl定义的 const NonceSize = 24
	var nonce [24]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	// This encrypts `plaintext` and appends the result to the nonce.
	encrypted := secretbox.Seal(nonce[:], plaintext, &nonce, &s.secretKey)
	return encrypted
}

// Decrypt 解密
func (s *NaCL) Decrypt(encrypted []byte) (plaintext []byte, err error) {
	if len(encrypted) <= 24 {
		return nil, fmt.Errorf("invalid ciphertext size (<=24)")
	}

	// When you decrypt, you must use the same nonce and key you used to
	// encrypt the message. One way to achieve this is to store the nonce
	// alongside the encrypted message. Above, we stored the nonce in the first
	// 24 bytes of the encrypted text.
	var decryptNonce [24]byte
	copy(decryptNonce[:], encrypted[:24])
	decrypted, ok := secretbox.Open(nil, encrypted[24:], &decryptNonce, &s.secretKey)
	if !ok {
		return nil, fmt.Errorf("decryption error")
	}

	return decrypted, nil
}
