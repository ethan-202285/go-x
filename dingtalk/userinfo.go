package dingtalk

// UserRole 钉钉用户角色
type UserRole struct {
	ID        uint   `json:"id"`        //
	Name      string ``                 // 角色名称
	GroupName string `json:"groupName"` // 角色组名称
}

// UserInfo 钉钉用户资料
type UserInfo struct {
	UnionID         string            `json:"unionid"`                     // "PiiiPyQqBNBii0HnCJ3zljcuAiEiE"，不会改变
	Remark          string            `json:"remark"`                      // "remark",
	UserID          string            `json:"userid" gorm:"primary_key"`   // "zhangsan"，创建后不可修改
	IsLeaderInDepts string            `json:"isLeaderInDepts"`             // "{1:false}",
	IsBoss          bool              `json:"isBoss"`                      // false,
	HiredDate       uint64            `json:"hiredDate"`                   // 1520265600000,
	IsSenior        bool              `json:"isSenior"`                    // false,
	Tel             string            `json:"tel"`                         // "xxx-xxxxxxxx", 分机号（仅限企业内部开发调用）
	Department      []int             `json:"department" gorm:"type:text"` // [1,2],
	WorkPlace       string            `json:"workPlace"`                   // "place",
	Email           string            `json:"email"`                       // "test@xxx.com",
	OrderInDepts    string            `json:"orderInDepts"`                // "{1:71738366882504}",
	Mobile          string            `json:"mobile"`                      // "1xxxxxxxxxx", 手机号码
	Active          bool              `json:"active"`                      // false,
	Avatar          string            `json:"avatar"`                      // "xxx",
	IsAdmin         bool              `json:"isAdmin"`                     // false, 是否为企业的管理员
	IsHide          bool              `json:"isHide"`                      // false,
	JobNumber       string            `json:"jobnumber"`                   // "001",
	Name            string            `json:"name"`                        // "张三",
	ExtAttr         map[string]string `json:"extattr" gorm:"type:jsonb"`   // {}, 扩展属性，可以设置多种属性
	StateCode       string            `json:"stateCode"`                   // "86",
	Position        string            `json:"position"`                    // "manager",
	Roles           []UserRole        `json:"roles" gorm:"type:jsonb"`     // [{"id": 149507744, "name": "总监", "groupName": "职务"}]
}
