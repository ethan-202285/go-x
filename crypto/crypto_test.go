package crypto

import (
	"encoding/base64"
	"fmt"
	"testing"
)

var (
	key         = []byte("12345678901234567890123456789012")
	content15   = []byte(`{"userID":112,}`)
	content31   = []byte(`{"userID":1,"name":"永安王"}`)
	content47   = []byte(`{"userID":12,"name":"永安王","sex":"female"}`)
	content63   = []byte(`{"userID":12,"name":"永安王","sex":"female","city":"Senzen"}`)
	content64   = []byte(`{"userID":12,"name":"永安王","sex":"female","city":"Senzhen"}`)
	contentLong = []byte(`{"userID":12,"name":"永安王","sex":"female","city":"Senzen", "userID":12,"name":"永安王","sex":"female","city":"Senzen", "userID":12,"name":"永安王","sex":"female","city":"Senzen"}`)
)

type Cryptor interface {
	Encrypt([]byte) []byte
	Decrypt([]byte) ([]byte, error)
}

func TestCBC(t *testing.T) {
	list := [][]byte{
		content15,
		content31,
		content47,
		content63,
		content64,
		contentLong,
	}

	cryptors := map[string]Cryptor{
		"CBC":  NewCBC(key),
		"GCM":  NewGCM(key),
		"NaCL": NewNaCL(key),
	}

	for _, plain := range list {
		fmt.Printf("原始:\t%s (%d字节)\n", plain, len(plain))
		for _, name := range []string{"CBC", "GCM", "NaCL"} {
			cryptor := cryptors[name]
			Run(t, name, cryptor, plain)
		}
		fmt.Printf("\n\n")
	}
}

func Run(t *testing.T, name string, cryptor Cryptor, plain []byte) {
	encrypted := cryptor.Encrypt(plain)
	encodeb64 := base64.RawURLEncoding.EncodeToString(encrypted)
	decrypted, err := cryptor.Decrypt(encrypted)
	if err != nil {
		panic(err)
	}
	rate := int(float32(len(plain)) / float32(len(encrypted)) * 100)
	fmt.Printf("[%s]\t%s (%d字节，base64后%d字符) 有效载荷：%d%%\n", name, encodeb64, len(encrypted), len(encodeb64), rate)
	fmt.Printf("\t%s (%d字节)\n", decrypted, len(decrypted))

	copy(encrypted[:1], []byte{0})
	decrypted, err = cryptor.Decrypt(encrypted)
	if err != nil {
		fmt.Printf("\t\033[32m无法篡改！\033[0m\n")
	} else {
		fmt.Printf("\033[31m篡改：\033[0m\t%s (%d字节)\n", decrypted, len(decrypted))
	}
}
