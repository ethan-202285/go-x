package dingtalk

// UserRole 钉钉用户角色
type UserRole struct {
	ID        uint   `json:"id"`        //
	Name      string `json:"name"`      // 角色名称
	GroupName string `json:"groupName"` // 角色组名称
}

// UserInfo 钉钉用户资料
type UserInfo struct {
	UnionID         string            `json:"unionid"`         // "PiiiPyQqBNBii0HnCJ3zljcuAiEiE"，不会改变
	Remark          string            `json:"remark"`          // "remark",
	UserID          string            `json:"userid"`          // "zhangsan"，创建后不可修改
	IsLeaderInDepts string            `json:"isLeaderInDepts"` // "{1:false}",
	IsBoss          bool              `json:"isBoss"`          // false,
	HiredDate       uint64            `json:"hiredDate"`       // 1520265600000,
	IsSenior        bool              `json:"isSenior"`        // false,
	Tel             string            `json:"tel"`             // "xxx-xxxxxxxx", 分机号（仅限企业内部开发调用）
	Department      []int             `json:"department"`      // [1,2],
	WorkPlace       string            `json:"workPlace"`       // "place",
	Email           string            `json:"email"`           // "test@xxx.com",
	OrderInDepts    string            `json:"orderInDepts"`    // "{1:71738366882504}",
	Mobile          string            `json:"mobile"`          // "1xxxxxxxxxx", 手机号码
	Active          bool              `json:"active"`          // false,
	Avatar          string            `json:"avatar"`          // "xxx",
	IsAdmin         bool              `json:"isAdmin"`         // false, 是否为企业的管理员
	IsHide          bool              `json:"isHide"`          // false,
	JobNumber       string            `json:"jobnumber"`       // "001",
	Name            string            `json:"name"`            // "张三",
	ExtAttr         map[string]string `json:"extattr"`         // {}, 扩展属性，可以设置多种属性
	StateCode       string            `json:"stateCode"`       // "86",
	Position        string            `json:"position"`        // "manager",
	Roles           []UserRole        `json:"roles"`           // [{"id": 149507744, "name": "总监", "groupName": "职务"}]
}

// Department 部门
type Department struct {
	ID                    int    `json:"id"`                    // 2,
	Name                  string `json:"name"`                  // "xxx",
	Order                 int    `json:"order"`                 // 10,
	ParentID              int    `json:"parentid"`              // 1,
	CreateDeptGroup       bool   `json:"createDeptGroup"`       // true,
	AutoAddUser           bool   `json:"autoAddUser"`           // true,
	DeptHiding            bool   `json:"deptHiding"`            // true,
	DeptPermits           string `json:"deptPermits"`           // "3|4",
	UserPermits           string `json:"userPermits"`           // "userid1|userid2",
	OuterDept             bool   `json:"outerDept"`             // true,
	OuterPermitDepts      string `json:"outerPermitDepts"`      // "1|2",
	OuterPermitUsers      string `json:"outerPermitUsers"`      // "userid3|userid4",
	OrgDeptOwner          string `json:"orgDeptOwner"`          // "manager1122",
	DeptManagerUseridList string `json:"deptManagerUseridList"` // "manager1122|manager3211",
	SourceIdentifier      string `json:"sourceIdentifier"`      // "source"
}
