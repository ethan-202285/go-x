package dingtalk

import (
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

	info = &UserInfo{}
	err = dd.Client.Get(r.Replace(url), info)
	if err != nil {
		return nil, err
	}
	return
}

// UserInfoByCode 根据免登授权码获取用户信息
func (dd *Dingtalk) UserInfoByCode(code string) (info *UserInfo, err error) {
	url := "https://oapi.dingtalk.com/user/getuserinfo?access_token=ACCESS_TOKEN&code=CODE"
	r := strings.NewReplacer("CODE", code)

	var result struct {
		UserID string `json:"userid"`
	}
	err = dd.Client.Get(r.Replace(url), &result)
	if err != nil {
		return nil, err
	}
	return dd.UserInfo(result.UserID)
}

// SendWorkMessage 发送消息(这个支持任何客服消息，但推荐用下面的快捷方法)
func (dd *Dingtalk) SendWorkMessage(message map[string]interface{}) (taskID uint64, err error) {
	message["agent_id"] = dd.config.AgentID
	url := "https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2?access_token=ACCESS_TOKEN"
	var result struct {
		TaskID uint64 `json:"task_id"`
	}
	err = dd.Client.PostJSON(url, message, &result)
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

	var result Department
	err := dd.Client.Get(r.Replace(url), &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// ListDepartment 列出所有子部门
func (dd *Dingtalk) ListDepartment(parentOpenID int) ([]*Department, error) {
	url := "https://oapi.dingtalk.com/department/list?access_token=ACCESS_TOKEN&fetch_child=true&id=PARENT_ID"
	r := strings.NewReplacer("PARENT_ID", strconv.Itoa(parentOpenID))

	var result struct {
		Department []*Department `json:"department"`
	}
	err := dd.Client.Get(r.Replace(url), &result)
	if err != nil {
		return nil, err
	}
	return result.Department, nil
}

// ListUserInDepartment 列出部门下用户详情
func (dd *Dingtalk) ListUserInDepartment(deptOpenID int) ([]*UserInfo, error) {
	url := "https://oapi.dingtalk.com/user/listbypage?access_token=ACCESS_TOKEN&department_id=DEPARTMENT_ID&offset=0&size=100"
	r := strings.NewReplacer("DEPARTMENT_ID", strconv.Itoa(deptOpenID))

	var result struct {
		Userlist []*UserInfo `json:"userlist"`
	}
	err := dd.Client.Get(r.Replace(url), &result)
	if err != nil {
		return nil, err
	}
	return result.Userlist, nil
}

type createWorkflowInstance struct {
	Request CreateWorkflowInstanceRequest `json:"request"`
}

// CreateWorkflowInstanceRequest 创建钉钉实例参数
type CreateWorkflowInstanceRequest struct {
	ProcessCode         string              `json:"process_code"`
	OriginatorUserID    string              `json:"originator_user_id"`
	Title               string              `json:"title"`
	FormComponentValues []WorkflowComponent `json:"form_component_values"`
	URL                 string              `json:"url"`
}

// WorkflowComponent 创建钉钉实例组件内容
type WorkflowComponent struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CreateWorkflowInstance 创建工作流实例
func (dd *Dingtalk) CreateWorkflowInstance(request CreateWorkflowInstanceRequest) (instanceID string, err error) {
	url := "https://oapi.dingtalk.com/topapi/process/workrecord/create?access_token=ACCESS_TOKEN"

	var payload = createWorkflowInstance{
		Request: request,
	}
	var result struct {
		Result struct {
			ProcessInstanceID string `json:"process_instance_id"`
		} `json:"result"`
	}
	err = dd.Client.PostJSON(url, payload, &result)
	if err != nil {
		return "", err
	}
	return result.Result.ProcessInstanceID, nil
}

// updateWorkflowInstance 更新钉钉实例参数
type updateWorkflowInstance struct {
	Request           UpdateWorkflowInstanceRequest `json:"request"`
	CancelRunningTask bool                          `json:"cancel_running_task"`
}

// UpdateWorkflowInstanceRequest 更新钉钉实例参数
type UpdateWorkflowInstanceRequest struct {
	ProcessInstanceID string `json:"process_instance_id"`
	Status            string `json:"status"`
	Result            string `json:"result"`
}

const (
	// WorkflowInstanceStatusCompleted 已完成
	WorkflowInstanceStatusCompleted = "COMPLETED"
	// WorkflowInstanceStatusTerminated 终止
	WorkflowInstanceStatusTerminated = "TERMINATED"
	// WorkflowInstanceResultAgree 同意
	WorkflowInstanceResultAgree = "agree"
	// WorkflowInstanceResultRefuse 拒绝
	WorkflowInstanceResultRefuse = "refuse"
)

// UpdateWorkflowInstance 更新工作流实例
func (dd *Dingtalk) UpdateWorkflowInstance(request UpdateWorkflowInstanceRequest, cancelRunningTask bool) error {
	url := "https://oapi.dingtalk.com/topapi/process/workrecord/update?access_token=ACCESS_TOKEN"

	var payload = updateWorkflowInstance{
		Request:           request,
		CancelRunningTask: cancelRunningTask,
	}
	return dd.Client.PostJSON(url, payload, nil)
}

type createWorkflowTask struct {
	Request CreateWorkflowTaskRequest `json:"request"`
}

// CreateWorkflowTaskRequest 待办事项
type CreateWorkflowTaskRequest struct {
	AgentID           int                      `json:"agentid"`
	ProcessInstanceID string                   `json:"process_instance_id"`
	ActivityID        string                   `json:"activity_id"`
	Tasks             []CreateWorkflowTaskNode `json:"tasks"`
}

// CreateWorkflowTaskNode 节点
type CreateWorkflowTaskNode struct {
	UserID string `json:"userid"`
	URL    string `json:"url"`
}

// CreateWorkflowTaskResponseNode 节点
type CreateWorkflowTaskResponseNode struct {
	UserID string `json:"userid"`
	TaskID int    `json:"task_id"`
}

// CreateWorkflowTask 创建钉钉待办事项
func (dd *Dingtalk) CreateWorkflowTask(request CreateWorkflowTaskRequest) ([]CreateWorkflowTaskResponseNode, error) {
	url := "https://oapi.dingtalk.com/topapi/process/workrecord/task/create?access_token=ACCESS_TOKEN"

	var payload = createWorkflowTask{
		Request: request,
	}
	var result struct {
		Tasks []CreateWorkflowTaskResponseNode `json:"tasks"`
	}
	err := dd.Client.PostJSON(url, payload, &result)
	if err != nil {
		return nil, err
	}
	return result.Tasks, nil
}

type updateWorkflowTask struct {
	Request UpdateWorkflowTaskRequest `json:"request"`
}

// UpdateWorkflowTaskRequest 待办事项
type UpdateWorkflowTaskRequest struct {
	AgentID           int                      `json:"agentid"`
	ProcessInstanceID string                   `json:"process_instance_id"`
	Tasks             []UpdateWorkflowTaskNode `json:"tasks"`
}

// UpdateWorkflowTaskNode 节点
type UpdateWorkflowTaskNode struct {
	TaskID int    `json:"task_id"`
	Status string `json:"status"`
	Result string `json:"result"`
}

// UpdateWorkflowTask 更新订单待办事项
func (dd *Dingtalk) UpdateWorkflowTask(request UpdateWorkflowTaskRequest) error {
	url := "https://oapi.dingtalk.com/topapi/process/workrecord/task/update?access_token=ACCESS_TOKEN"

	var payload = updateWorkflowTask{
		Request: request,
	}
	return dd.Client.PostJSON(url, payload, nil)
}
