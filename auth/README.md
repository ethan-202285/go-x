
auth 模块
=============
提供user、token、identity的数据模型  
提供数据操作的Repository  
提供Login、Renew、Logout的Service  
提供从r.Context()里面设置/读取user的方法  
提供Authenticated的middleware  
提供Login、Renew、Logout的Handler，以及一个DemoPage  


用法
-------------
* 实例化
    ```go
    import (
        "github.com/goodwong/go-x/auth"
    )

    dsn := fmt.Sprintf("host=db port=5432 user=app dbname=app password=app sslmode=disable")
    db, err := gorm.Open("postgres", dsn)
    secretKey := []byte("aasdfkjksjdfaaasdfkjksjdfa123405") // 32 bytes
    auths := auth.New(auth.Config{DB: db, SecretKey: secretKey})

    // 如果需要创建数据库
    auths.Repository.AutoMigrate()
    ```

* 添加密码登陆方式
    ```go
    import (
        "github.com/goodwong/go-x/auth/providers/password"
    )

	passwords := password.NewProvider(&password.Config{
		Auth:      auths,
		SecretKey: secretKey,
	})
	auths.RegisterProvider(passwords)

    // passwords 可能用于创建密码
    if _, err := passwords.Register(username, password, userIDs...); err != nil {
        log.Fatal(err)
    }
    // 或重置密码
    if err := passwords.SetPassword(username, password); err != nil {
        log.Fatal(err)
    }
    ```

* 添加到路由规则
    ```go
    r := chi.NewRouter()
    r.Handle("/api/login", auths.Handler.Mux())

    // 或者分别路由
    r.Route("/api/login", func(r chi.Router) {
        // Login Demo Page
        r.Get("/", auths.Handler.LoginDemoPage)
        // Login
        r.Post("/", auths.Handler.HandleLogin)
        // Renew
        r.Patch("/", auths.Handler.HandleRenew)
        // Logout
        authenticated := auths.Middleware.AuthenticatedWithUser
        r.With(authenticated).Delete("/", auths.Handler.HandleLogout)
    })
    ```

    > 前端开发人员，可以访问`Login Demo Page`，查看登陆演示，方便理解


功能
-------------
* 从request.Contex()获取用户ID
    ```go
    ctx := auth.NewContext(req.Context())
    // 设置用户ID
    ctx.WithUserID(125)
    
    // 获取
    userID := ctx.UserID()
    if userID == 0 {
        log.Fatal("invalid userID")
    }
	user := ctx.User() // return &User{ID:125}
	if user == nil {
		log.Fatal("invalid user")
	}
    
    // 设置用户
    ctx.WithUser(&User{ID: 15})
    
    // 用在middleware中传递下去
    func (m *Middleware) Authenticated(next http.Handler) http.Handler {
        attach := func(next http.Handler) http.Handler {
            return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                ctx := auth.NewContext(r.Context())

                // load
                user, _ := auths.Repository.Find(ctx.UserID())
                ctx.WithUser(user)

                // 传递下去
                next.ServeHTTP(w, ctx.AttachRequest(r))
            })
        }
	    return attach(next)
    }
    ```

* 中间件
    ```go
    // 确保用户已登录，并在context里设置userID
    // 鉴权方式：
    //   1. jwt
    //   2. cookie中的refresh_token
    r := chi.NewRouter()
    authenticated := auths.Middleware.Authenticated
    logout := auths.Handler.HandleLogout
    r.With(authenticated).Delete("/", logout)

    // 在前面的基础上，每次还会从数据库中拉取user信息
    auths.Middleware.AuthenticatedWithUser
    ```

高级
-------------

* 登陆注销功能
    > 大多数情况下使用handler和middleware就可以了，  
    > 当然你愿意的话，自己写handler也是可以的
    ```go
    // 登陆
    tokens, err := auths.Service.Login(provider, payload, remember, device)
	setCookie(w, "jwt", tokens.Token, tokens.TokenExpires)
	if tokens.RefreshToken != nil {
		setCookie(w, "refresh_token", *tokens.RefreshToken, *tokens.RefreshTokenExpires)
    }
    
    // token续约
    user, tokens, err := h.auth.Service.Renew(params.RefreshToken)
    setCookie(w, "jwt", tokens.Token, tokens.TokenExpires)

    // 登出
    auths.Service.Logout(user, device)
	deleteCookie(w, "jwt")
	deleteCookie(w, "refresh_token")
    ```
* 检查jwt是否提前失效（用户主动注销）
    > 常规的jwt是无法主动失效的  
    > 但是我们提供了方法，能够检测到用户主动注销掉的jwt
    ```go
    needRelogin := auth.Service.JwtInvalid(jwtToken)
    ```


* 添加自定义provider
    > 
    ```go
    // 1. 实现`LoginProvider`接口
    type LoginProvider interface {
        Name() string // 返回provider的名字
        Login(credentials []byte) (user *User, err error)
    }

    // 2. 注册登陆方式
    auths.RegisterProvider(passwords)
    ```