package weapp_test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"log"
	"os"
	"testing"

	"github.com/goodwong/go-x/wechat/weapp"
)

var (
	wechatWeapp *weapp.Weapp
	code        string
)

// 测试之前，请设置Env变量，命令行：
// export CGO_ENABLED=0
// AppID= AppSecret=
// go test ./weapp
// 或
// AppID= AppSecret= go test ./weapp
func init() {
	config := weapp.Config{
		AppID:     os.Getenv("AppID"),
		AppSecret: os.Getenv("AppSecret"),
	}
	code = os.Getenv("code")
	if code == "" {
		log.Print("\033[31m缺少参数code，跳过\033[0m\n" +
			"\t\033[33m请在小程序wx.login()获取code，然后按如下格式运行测试：\n" +
			"\t\033[7mcode= \033[0;33m go test ./wechat/weapp\033[0m",
		)
	}
	wechatWeapp = weapp.New(&config)
}

func Test_GetAccessToken(t *testing.T) {
	tokenString, err := wechatWeapp.Client.GetAccessToken()
	if err != nil || len(tokenString) == 0 {
		t.Fatal("GetAccessToken失败：", err)
	}
	t.Log("获取AccessToken成功", tokenString)
}

func TestCode2session(t *testing.T) {
	if code == "" {
		return
	}
	result, err := wechatWeapp.Code2session(code)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v\n", result)
}

func TestUnlimitedWacode(t *testing.T) {
	// respBytes, err := wechatWeapp.UnlimitedWacode("pages/gift/index/main", "gift=123")
	// if err != nil {
	// 	t.Logf("生成小程序码需要已经发布的页面，测试号无法测试")
	// 	t.Fatal(err)
	// }
	//
	// output, err := os.Create("create_unlimited_wxcode_output.png")
	// if err != nil {
	// 	t.Fatal("创建output文件失败：", err)
	// }
	// defer output.Close()
	// output.Write(respBytes)
}

func TestDecrpyt(t *testing.T) {

	// 小程序返回数据：
	// {"errMsg":"getUserInfo:ok","rawData":"{\"nickName\":\"老小王\",\"gender\":1,\"language\":\"zh_CN\",\"city\":\"Shenzhen\",\"province\":\"Guangdong\",\"country\":\"China\",\"avatarUrl\":\"https://wx.qlogo.cn/mmopen/vi_32/icX0N3ibGic4Vx0ZhGxNvUictmw3cKhrrKsrtzFYpvIGcbqBbCDM3PGkibyvEkunaC1M1jnQxPWh7eufkZtgbBCmuqw/132\"}","userInfo":{"nickName":"老小王","gender":1,"language":"zh_CN","city":"Shenzhen","province":"Guangdong","country":"China","avatarUrl":"https://wx.qlogo.cn/mmopen/vi_32/icX0N3ibGic4Vx0ZhGxNvUictmw3cKhrrKsrtzFYpvIGcbqBbCDM3PGkibyvEkunaC1M1jnQxPWh7eufkZtgbBCmuqw/132"},"signature":"8ef5fafde244b94fbbfd2e7e2b5c457ab253aeb3","encryptedData":"W5ptNJ/HWKYgWtS55M9xf1N5HlLfCZ9jTIhoTx59e749j55i/+Klfe8i3uZ5af1s/pXChszniNLncr6ypqod/6Ju1HmbsAh2bsIYpgkWgF+pLqxNqXOkVTHyphw23yriQkNqT18042n19F8oZ8aR/jj1Q04vP/TQJyQvhxAjpTJREyMVsUGBc+CNiId5vZDWu+vWWuL9+W2AMLbJEJ2zojNZ1OL8nx2JzVl1PBmJwaVgVfjFtyc9F0prgD2H14QSvn4ID66m5B7eIbLB4CUv5X5eAgy9oRQZGALSI8f8U3krjR0B/FpAOr/Hb3SBGpRHec7Z3W36uCnXhKnQFvTmRIN3FTL2TK1JagoEr4xQNjOEGKvh6zPhnomz51TZgeSjcI0TudmyQmuSyyL8MGVX/ScSOsCyNaLy03NK6YcrIkqkdbobx7R6E0OMd0F7rhniZTBUHqm7/iNYmUWa2ME42iKXoCtnilgR2N/wpIFOGfc=","iv":"ysEVDxaTYM4dgXavukQywA=="}

	// 解密
	// ------------------------------------------->Z<------
	ciphertextB64 := "W5ptNJ/HWKYgWtS55M9xf1N5HlLfCZ9jTIhoTx59e749j55i/+Klfe8i3uZ5af1s/pXChszniNLncr6ypqod/6Ju1HmbsAh2bsIYpgkWgF+pLqxNqXOkVTHyphw23yriQkNqT18042n19F8oZ8aR/jj1Q04vP/TQJyQvhxAjpTJREyMVsUGBc+CNiId5vZDWu+vWWuL9+W2AMLbJEJ2zojNZ1OL8nx2JzVl1PBmJwaVgVfjFtyc9F0prgD2H14QSvn4ID66m5B7eIbLB4CUv5X5eAgy9oRQZGALSI8f8U3krjR0B/FpAOr/Hb3SBGpRHec7Z3W36uCnXhKnQFvTmRIN3FTL2TK1JagoEr4xQNjOEGKvh6zPhnomz51TZgeSjcI0TudmyQmuSyyL8MGVX/ScSOsCyNaLy03NK6YcrIkqkdbobx7R6E0OMd0F7rhniZTBUHqm7/iNYmUWa2ME42iKXoCtnilgR2N/wpIFOGfc="
	ivB64 := "ysEVDxaTYM4dgXavukQywA=="
	sessionKey := "trOV+Gdz2vjsQIKPtpyjtw=="

	ciphertext, err := base64.StdEncoding.DecodeString(ciphertextB64)
	if err != nil {
		t.Fatal(err)
	}
	iv, err := base64.StdEncoding.DecodeString(ivB64)
	if err != nil {
		t.Fatal(err)
	}
	secretKey, err := base64.StdEncoding.DecodeString(sessionKey)
	if err != nil {
		t.Fatal(err)
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		t.Fatal(err)
	}
	if len(ciphertext) <= aes.BlockSize {
		t.Fatal(err)
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		t.Fatal(err)
	}
	plaintext := make([]byte, len(ciphertext))
	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(plaintext, ciphertext)

	plaintext = pkcs7UnPadding(plaintext)
	t.Log(string(plaintext))

	// !小程序的加密有漏洞
	// AES CBC没有校验功能，而小程序的 signature 不负责encryptedData部分的内容进行校验
	// 也就是encryptedData可以修改，后端无从检查发现。

	// 校验
	rawData := `{"nickName":"老小王","gender":1,"language":"zh_CN","city":"Shenzhen","province":"Guangdong","country":"China","avatarUrl":"https://wx.qlogo.cn/mmopen/vi_32/icX0N3ibGic4Vx0ZhGxNvUictmw3cKhrrKsrtzFYpvIGcbqBbCDM3PGkibyvEkunaC1M1jnQxPWh7eufkZtgbBCmuqw/132"}`
	// signature:"8ef5fafde244b94fbbfd2e7e2b5c457ab253aeb3"
	h := sha1.New()
	h.Write([]byte(rawData + sessionKey))
	t.Logf("%x\n", string(h.Sum(nil)))
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
