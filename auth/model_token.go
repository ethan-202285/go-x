package auth

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"io"
	"time"

	"github.com/goodwong/go-x/crypto"

	"golang.org/x/crypto/bcrypt"
)

// DefaultRefreshTokenLife 默认token有效时长
var DefaultRefreshTokenLife = 365 * 24 * time.Hour // 1 Year

// Usage:
// token := newToken(Token{UserID:13, Device:"web", Remark:"浏览器"}, KEY)
func newToken(params Token, secretKey []byte) *Token {
	if params.UserID == 0 {
		panic("params.UserID未设置")
	}
	if len(params.Device) == 0 {
		panic("params.Device未设置")
	}
	now := time.Now()
	t := &Token{
		UserID: params.UserID,
		Device: params.Device,
		Remark: params.Remark,
	}
	t.IssuedAt = now
	t.ExpiredAt = now.Add(DefaultRefreshTokenLife)
	// 生成12位随机数
	// bcrypt调试到最低1ms左右
	// 因为是按照字节随机的无规律，充分利用每字节128种可见字符可能性，无法做生日攻击
	// 所以12位随机是: 2**(7*12) = 3.8E28
	// 假设每次1ms，则1s内可以运行1000次
	// 破解时间： 2**(8*12)/(1000*86400*365) = 1.2E18 年
	t.secretKey = secretKey
	t.nonce = make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, t.nonce); err != nil {
		panic(err.Error())
	}
	// sign
	t.Hash = t.sign()

	return t
}

func parseTokenString(tokenString string, secretKey []byte) (t *Token, err error) {
	// base64 解码
	data, err := base64.RawURLEncoding.DecodeString(tokenString)
	if err != nil {
		return nil, ErrInvalidToken
	}

	// 解密
	plaintext, err := crypto.NewNaCL(secretKey).Decrypt(data)
	if err != nil {
		return nil, err
	}

	// 解包
	ID := uint64(binary.BigEndian.Uint64(plaintext[:8]))
	nonce := plaintext[8 : 8+12]
	if ID == 0 {
		return nil, ErrInvalidToken
	}

	// 剩下的更多内容需要去数据库获取
	return &Token{
		ID:        ID,
		nonce:     []byte(nonce),
		secretKey: secretKey,
	}, nil
}

// Token 登录后的token （一个用户可以由多个token）
type Token struct {
	ID        uint64
	UserID    uint64
	Device    string // 比如“home”“office”，有前端程序定义
	Remark    string // 比如“家”、“办公室”，用户定义
	Hash      string //
	IssuedAt  time.Time
	ExpiredAt time.Time
	DeletedAt *time.Time
	// 临时变量
	secretKey []byte
	nonce     []byte
}

// TableName 指定数据表名(gorm)
func (t *Token) TableName() string {
	return "user_tokens"
}

// ExpiresAfter 设置token有效期，很少用到，所以独立方法
func (t *Token) ExpiresAfter(duration time.Duration) *Token {
	t.ExpiredAt = t.IssuedAt.Add(duration)
	return t
}

// sign 生成签名
// 由newToken()自动调用
func (t *Token) sign() string {
	if len(t.nonce) == 0 {
		panic("缺少Token.nonce，只有newToken()创建的对象，才能导出Signature")
	}
	if len(t.secretKey) == 0 {
		panic("缺少Token.secretKey，只有newToken()创建的对象，才能导出Signature")
	}

	// 打包
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, t.nonce)
	// 加这么多参数没有意义，反而破坏了密码的复杂性和不可预知性
	//binary.Write(buf, binary.BigEndian, t.ID)
	//binary.Write(buf, binary.BigEndian, t.UserID)
	//binary.Write(buf, binary.BigEndian, []byte(t.Device))
	//binary.Write(buf, binary.BigEndian, uint64(t.IssuedAt.UTC().Unix()))
	//binary.Write(buf, binary.BigEndian, t.secretKey)

	// 其实这里用sha1完全可以
	// 因为不是用户密码，字典攻击无效，穷举费时
	// （P.S. 如果是用户密码，如果是弱口令，bcrypt也救不了）
	// bcryptz自动加盐，更加难做暴力破解
	//
	// sha2
	//h := sha256.New()
	//h.Write(buf.Bytes())
	//hash := h.Sum(nil)
	//
	// bcrypt
	hash, err := bcrypt.GenerateFromPassword(buf.Bytes(), bcrypt.DefaultCost) // 50ms
	if err != nil {
		panic(err)
	}

	// base64
	//return base64.URLEncoding.EncodeToString(hash)
	return string(hash) // bcrypt 本身有base64编码
}

// Verify 验证
func (t *Token) Verify() bool {
	if len(t.Hash) == 0 {
		panic("缺少Token.Hash，请查询数据库获取Hash，再Verify()")
	}

	// 打包
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, t.nonce)
	//binary.Write(buf, binary.BigEndian, t.secretKey)

	err := bcrypt.CompareHashAndPassword([]byte(t.Hash), buf.Bytes())
	return err == nil
}

// TokenString 导出AES加密Base64编码后的ID+nonce字符串
func (t *Token) TokenString() string {
	if t.ID == 0 {
		panic("缺少Token.ID，必须保存到数据库，才能导出String")
	}
	if len(t.nonce) == 0 {
		panic("缺少Token.nonce，只有newToken()创建的对象，才能导出String")
	}
	if len(t.secretKey) == 0 {
		panic("缺少Token.secretKey，只有newToken()创建的对象，才能导出String")
	}

	// 打包
	// |-id(8)-|-nonce(12)-| (20字节) 新方案
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, t.ID)    // 8 bytes
	binary.Write(buf, binary.BigEndian, t.nonce) // 12 bytes

	// 加密
	encrypted := crypto.NewNaCL(t.secretKey).Encrypt(buf.Bytes())

	// base64编码
	return base64.RawURLEncoding.EncodeToString(encrypted)
}
