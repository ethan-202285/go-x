package dingtalk

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"

	"github.com/goodwong/go-x/dingtalk/client"
)

// New 创建
func New(cfg *Config) *Dingtalk {
	r := strings.NewReplacer("KEY", cfg.AppKey, "SECRET", cfg.AppSecret)
	tokenAPI := r.Replace("https://oapi.dingtalk.com/gettoken?appkey=KEY&appsecret=SECRET")

	client := client.New(&client.Config{
		TokenAPI: tokenAPI,
	})
	return &Dingtalk{
		config: cfg,
		Client: client,
	}
}

// Config 配置类
type Config struct {
	CorpID    string
	AgentID   uint64
	AppKey    string
	AppSecret string
}

// Dingtalk 功能类
type Dingtalk struct {
	config *Config
	Client *client.APIClient
}

// UserInfo 用户信息
// 如果数据库有信息，直接返回
// 更新用户建议用钉钉的主动通知接口，效率更高，体验更好
func (dd *Dingtalk) UserInfo(userID string) (info *UserInfo, err error) {
	// 从服务器拉取信息
	url := "https://oapi.dingtalk.com/user/get?access_token=ACCESS_TOKEN&userid=USERID"
	r := strings.NewReplacer("USERID", userID)

	respBytes, err := dd.Client.Get(r.Replace(url))
	if err != nil {
		return nil, err
	}
	info = &UserInfo{}
	err = json.Unmarshal(respBytes, &info)
	if err != nil {
		log.Println("json.Unmarshal失败，钉钉个人信息jsonString：", string(respBytes))
		return nil, err
	}
	return
}

// UserInfoByCode 根据免登授权码获取用户信息
func (dd *Dingtalk) UserInfoByCode(code string) (info *UserInfo, err error) {
	url := "https://oapi.dingtalk.com/user/getuserinfo?access_token=ACCESS_TOKEN&code=CODE"
	r := strings.NewReplacer("CODE", code)

	respBytes, err := dd.Client.Get(r.Replace(url))
	if err != nil {
		return nil, err
	}
	var result struct {
		UserID string `json:"userid"`
	}
	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		return nil, err
	}
	return dd.UserInfo(result.UserID)
}

// SendWorkMessage 发送消息(这个支持任何客服消息，但推荐用下面的快捷方法)
func (dd *Dingtalk) SendWorkMessage(message map[string]interface{}) (taskID uint64, err error) {
	message["agent_id"] = dd.config.AgentID
	url := "https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2?access_token=ACCESS_TOKEN"
	respBytes, err := dd.Client.PostJSON(url, message)
	if err != nil {
		return 0, err
	}
	var result struct {
		TaskID uint64 `json:"task_id"`
	}
	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		return 0, err
	}
	return result.TaskID, nil
}

// SendText 发送文字
func (dd *Dingtalk) SendText(receiver map[string]interface{}, content string) (taskID uint64, err error) {
	data := map[string]interface{}{
		"msg": map[string]interface{}{
			"msgtype": "text",
			"text": map[string]string{
				"content": content,
			},
		},
	}
	if userIDList := receiver["userid_list"]; userIDList != nil {
		data["userid_list"] = userIDList
	} else if deptIDList := receiver["dept_id_list"]; deptIDList != nil {
		data["dept_id_list"] = deptIDList
	} else if toAllUser := receiver["to_all_user"]; toAllUser != nil {
		data["to_all_user"] = true
	}
	return dd.SendWorkMessage(data)
}

// GetDepartment 获取部门详情
func (dd *Dingtalk) GetDepartment(deptOpenID int) (*Department, error) {
	url := "https://oapi.dingtalk.com/department/get?access_token=ACCESS_TOKEN&id=DEPARTMENT_ID"
	r := strings.NewReplacer("DEPARTMENT_ID", strconv.Itoa(deptOpenID))

	respBytes, err := dd.Client.Get(r.Replace(url))
	if err != nil {
		return nil, err
	}
	var result Department
	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListDepartment 列出所有子部门
func (dd *Dingtalk) ListDepartment(parentOpenID int) ([]*Department, error) {
	url := "https://oapi.dingtalk.com/department/list?access_token=ACCESS_TOKEN&fetch_child=true&id=PARENT_ID"
	r := strings.NewReplacer("PARENT_ID", strconv.Itoa(parentOpenID))

	respBytes, err := dd.Client.Get(r.Replace(url))
	if err != nil {
		return nil, err
	}
	var result struct {
		Department []*Department `json:"department"`
	}
	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		log.Printf("钉钉返回：%s\n", respBytes)
		return nil, err
	}
	return result.Department, nil
}

// ListUserInDepartment 列出部门下用户详情
func (dd *Dingtalk) ListUserInDepartment(deptOpenID int) ([]*UserInfo, error) {
	url := "https://oapi.dingtalk.com/user/listbypage?access_token=ACCESS_TOKEN&department_id=DEPARTMENT_ID&offset=0&size=100"
	r := strings.NewReplacer("DEPARTMENT_ID", strconv.Itoa(deptOpenID))

	respBytes, err := dd.Client.Get(r.Replace(url))
	if err != nil {
		return nil, err
	}
	var result struct {
		Userlist []*UserInfo `json:"userlist"`
	}
	err = json.Unmarshal(respBytes, &result)
	if err != nil {
		return nil, err
	}
	return result.Userlist, nil
}
