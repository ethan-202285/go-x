package auth

// LoginProvider 登陆方法
type LoginProvider interface {
	Name() string
	Login(credentials []byte) (user *User, err error)
}

// RegisterProvider 注册登陆方式
func (auth *Auth) RegisterProvider(provider LoginProvider) {
	name := provider.Name()
	auth.Service.providers[name] = provider
	return
}
