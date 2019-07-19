package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
)

// NewCBC 新建
func NewCBC(secretKey []byte) *CBC {
	return &CBC{secretKey: secretKey}
}

// CBC CBC加密
type CBC struct {
	secretKey []byte
}

// Encrypt 加密
func (a *CBC) Encrypt(content []byte) []byte {
	plaintext := pkcs7Padding(content, aes.BlockSize)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))

	// 随机iv
	if _, err := io.ReadFull(rand.Reader, ciphertext[:aes.BlockSize]); err != nil {
		panic(err)
	}
	iv := ciphertext[:aes.BlockSize]

	// AES CBC加密，IV放在前缀
	block, err := aes.NewCipher(a.secretKey)
	if err != nil {
		panic(err)
	}
	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(ciphertext[aes.BlockSize:], plaintext)
	return ciphertext
}

// Decrypt 解密
func (a *CBC) Decrypt(ciphertext []byte) (plaintext []byte, err error) {
	if len(ciphertext)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("ciphertext is not a multiple of the block size")
	}
	if len(ciphertext) <= aes.BlockSize {
		return nil, fmt.Errorf("invalid ciphertext size (<=%d)", aes.BlockSize)
	}

	block, err := aes.NewCipher(a.secretKey)
	if err != nil {
		return nil, err
	}
	plaintext = make([]byte, len(ciphertext)-aes.BlockSize)
	iv := ciphertext[:aes.BlockSize]
	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(plaintext, ciphertext[aes.BlockSize:])

	plaintext = pkcs7UnPadding(plaintext)
	return
}

func pkcs7Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize        //需要padding的数目
	padtext := bytes.Repeat([]byte{byte(padding)}, padding) //生成填充的文本
	return append(ciphertext, padtext...)
}

func pkcs7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func zeroPadding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{0}, padding) //用0去填充
	return append(ciphertext, padtext...)
}

func zeroUnPadding(origData []byte) []byte {
	return bytes.TrimFunc(origData,
		func(r rune) bool {
			return r == rune(0)
		})
}
