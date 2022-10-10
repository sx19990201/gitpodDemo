package domain

import (
	"context"
	"encoding/json"
	"github.com/vektah/gqlparser/v2/ast"
	"net/http"
)

const (
	Queries       = "queries"
	Mutations     = "mutations"
	Subscriptions = "subscriptions"

	query    = "query"
	mutation = "mutation"

	Legitimate   = 1 // 合法
	UnLegitimate = 2 // 非法

	Public                = 1
	Private               = 2
	GenerateSchemaGraphql = "wundergraph.app.schema.graphql"
)

// operations 指令
const (
	// InternalOperation 内部
	InternalOperation = "internalOperation"
	RBAC              = "rbac"
)

// operations 用户规则
const (
	RequireMatchAll = "requireMatchAll"
	RequireMatchAny = "requireMatchAny"
	DenyMatchAll    = "denyMatchAll"
	DenyMatchAny    = "denyMatchAny"
)

// FbOperations 操作api
type FbOperations struct {
	ID            int64  `db:"id" json:"id"`
	Method        string `db:"method" json:"method"`                // 方法类型 GET、POST、PUT、DELETE
	OperationType string `db:"operation_type" json:"operationType"` // 请求类型 queries,mutations,subscriptions
	IsPublic      int64  `db:"is_public" json:"isPublic"`           // 状态 1公有 2私有
	Remark        string `db:"remark" json:"remark"`                // 说明
	Legal         int64  `db:"legal" json:"legal"`                  // 是否合法 1合法 2非法
	Path          string `db:"path" json:"title"`                   // 路径
	Content       string `db:"content" json:"content"`              // 内容
	Enable        int64  `db:"enable" json:"enable"`                // 开关 true开 false关
	CreateTime    string `db:"create_time" json:"createTime"`
	UpdateTime    string `db:"update_time" json:"updateTime"`
	IsDel         int64  `db:"is_del" json:"isDel"`
	RoleType      string `db:"role_type" json:"roleType"`
	Roles         string `db:"roles" json:"roles"`
}

func (f *FbOperations) SetField(schemaDocument *ast.QueryDocument) {
	if schemaDocument == nil || len(schemaDocument.Operations) == 0 {
		return
	}
	// 解析成功获取查询类型
	if schemaDocument.Operations[0].Operation == mutation {
		f.Method = http.MethodPost
		f.OperationType = Mutations
	}
	if schemaDocument.Operations[0].Operation == query {
		f.Method = http.MethodGet
		f.OperationType = Queries
	}
	// TODO 订阅
	//if schemaDocument.Operations[0].Operation == query {
	//	f.Method = http.MethodGet
	//	f.OperationType = Queries
	//}

	for _, directive := range schemaDocument.Operations[0].Directives {
		// 获取公开状态
		if directive.Name == InternalOperation {
			f.IsPublic = Private
			continue
		}
		// 获取用户角色权限
		if directive.Name == RBAC {
			if len(directive.Arguments) == 0 {
				continue
			}
			rbacDocument := directive.Arguments[0]
			f.RoleType = rbacDocument.Name
			roleArr := make([]string, 0)
			if rbacDocument.Value == nil {
				continue
			}
			for _, row := range rbacDocument.Value.Children {
				roleArr = append(roleArr, row.Value.Raw)
			}
			roleBytes, _ := json.Marshal(roleArr)
			f.Roles = string(roleBytes)
			continue
		}
	}
}

type FbOperationsResult struct {
	ID            int64  `db:"id" json:"id"`
	Method        string `db:"method" json:"method"`                // 方法类型 GET、POST、PUT、DELETE
	OperationType string `db:"operation_type" json:"operationType"` // 请求类型 queries,mutations,subscriptions
	IsPublic      bool   `db:"is_public" json:"isPublic"`           // 状态 1公有 2私有
	Remark        string `db:"remark" json:"remark"`                // 说明
	Legal         bool   `db:"legal" json:"legal"`                  // 是否合法 1合法 2非法
	Path          string `db:"path" json:"path"`                    // 路径
	IsDir         bool   `db:"-" json:"isDir"`                      // 是否是文件夹
	Content       string `db:"content" json:"content"`              // 内容
	Enable        bool   `db:"enable" json:"enable"`                // 开关 true开 false关
	CreateTime    string `db:"create_time" json:"createTime"`
	UpdateTime    string `db:"update_time" json:"updateTime"`
	IsDel         int64  `db:"is_del" json:"-"`
}

func (f *FbOperationsResult) Transform() FbOperations {
	var legal int64 = 0
	var isPublic int64 = 0
	var enable int64 = 0

	if !f.Legal {
		legal = 1
	}
	if !f.IsPublic {
		isPublic = 1
	}
	if !f.Enable {
		enable = 1
	}
	return FbOperations{
		ID:            f.ID,
		Method:        f.Method,
		OperationType: f.OperationType,
		IsPublic:      isPublic,
		Remark:        f.Remark,
		Legal:         legal,
		Path:          f.Path,
		Content:       f.Content,
		Enable:        enable,
		CreateTime:    f.CreateTime,
		UpdateTime:    f.UpdateTime,
		IsDel:         f.IsDel,
	}

}

func (f *FbOperations) TransformToResult() FbOperationsResult {
	legal := true
	isPublic := true
	enable := true

	if f.Legal == 1 {
		legal = false
	}
	if f.IsPublic == 1 {
		isPublic = false
	}
	if f.Enable == 1 {
		enable = false
	}
	return FbOperationsResult{
		ID:            f.ID,
		Method:        f.Method,
		OperationType: f.OperationType,
		IsPublic:      isPublic,
		Remark:        f.Remark,
		Legal:         legal,
		Path:          f.Path,
		Content:       f.Content,
		Enable:        enable,
		CreateTime:    f.CreateTime,
		UpdateTime:    f.UpdateTime,
		IsDel:         f.IsDel,
	}

}

type OperationsUseCase interface {
	Store(ctx context.Context, f *FbOperations) (int64, error)
	ReName(ctx context.Context, oldName, newName string, enable int64) error
	ReDicName(ctx context.Context, oldName, newName string) error
	GetByPath(ctx context.Context, path string) (FbOperations, error)
	GetByID(ctx context.Context, id int64) (FbOperations, error)
	Update(ctx context.Context, f *FbOperations) (int64, error)
	Delete(ctx context.Context, id int64) (int64, error)
	FindOperations(ctx context.Context) ([]FbOperations, error)
	FindByPath(ctx context.Context, path string) ([]FbOperations, error)
	BatchDelete(c context.Context, opers []FbOperations) error
}

type OperationsRepository interface {
	Store(ctx context.Context, f *FbOperations) (int64, error)
	Update(ctx context.Context, f *FbOperations) (int64, error)
	GetByPath(ctx context.Context, path string) (FbOperations, error)
	GetByID(ctx context.Context, id int64) (FbOperations, error)
	Delete(ctx context.Context, id int64) (int64, error)
	FindOperations(ctx context.Context) ([]FbOperations, error)
	FindByPath(ctx context.Context, path string) ([]FbOperations, error)
	ReName(ctx context.Context, oldName, newName string, enable int64) error
	BatchDelete(c context.Context, ids string) (err error)
}
