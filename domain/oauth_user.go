package domain

import (
	"encoding/json"
	"fmt"
	"github.com/prisma/prisma-client-go/generator/ast/dmmf"

	//"github.com/prisma/prisma-client-go/generator/ast/dmmf"
	"strings"
)

type OauthUser struct {
	ID                 string `db:"id" json:"id"`                                   // id
	Name               string `db:"name" json:"name"`                               // 姓名
	NickName           string `db:"nick_name" json:"nick_name"`                     // 昵称
	UserName           string `db:"user_name" json:"user_name"`                     // 用户名
	EncryptionPassword string `db:"encryption_password" json:"encryption_password"` // 加密后密码
	Mobile             string `db:"mobile" json:"mobile"`                           // 手机号
	Email              string `db:"email" json:"email"`                             // 邮箱
	LastLoginTime      string `db:"last_login_time" json:"last_login_time"`         // 最后一次登陆时间
	Status             int64  `db:"status" json:"status"`                           // 状态
	MateData           string `db:"mate_data" json:"mate_data"`                     // 其他信息(json字符串保存)
	CreateTime         string `db:"create_time" json:"create_time"`                 // 创建时间
	UpdateTime         string `db:"update_time" json:"update_time"`                 // 修改时间
	IsDel              int64  `db:"is_del"  json:"is_del"`                          // 是否删除
}

type OauthUserResp struct {
	ID                 string `json:"id"`                 // id
	Name               string `json:"name"`               // 姓名
	NickName           string `json:"nickName"`           // 昵称
	UserName           string `json:"userName"`           // 用户名
	EncryptionPassword string `json:"encryptionPassword"` // 加密后密码
	Mobile             string `json:"mobile"`             // 手机号
	Email              string `json:"email"`              // 邮箱
	LastLoginTime      string `json:"lastLoginTime"`      // 最后一次登陆时间
	Status             int64  `json:"status"`             // 状态
	MateData           string `json:"mateData"`           // 其他信息(json字符串保存)
	CreateTime         string `json:"createTime"`         // 创建时间
	UpdateTime         string `json:"updateTime"`         // 修改时间
	IsDel              int64  `json:"is_del"`             // 是否删除
}

func (o *OauthUser) TransformResp() OauthUserResp {
	return OauthUserResp{
		ID:                 o.ID,
		Name:               o.Name,
		NickName:           o.NickName,
		UserName:           o.UserName,
		EncryptionPassword: o.EncryptionPassword,
		Mobile:             o.Mobile,
		Email:              o.Email,
		LastLoginTime:      o.LastLoginTime,
		Status:             o.Status,
		MateData:           o.MateData,
		CreateTime:         o.CreateTime,
		UpdateTime:         o.UpdateTime,
		IsDel:              o.IsDel,
	}
}

type MateData struct {
	CountryCode   string   `json:"countryCode"`
	Country       string   `json:"country"`
	City          string   `json:"city"`
	Province      string   `json:"province"`
	StreetAddress string   `json:"streetAddress"`
	ExternalId    string   `json:"externalId"`
	PostalCode    string   `json:"postalCode"`
	SendAddress   bool     `json:"sendAddress"`
	Role          []FbRole `json:"role"`
	ProfileClaims
	EmailClaims
	AddressClaims
	PhoneClaims
}

const (
	ProfileScope = "profile"
	EmailScope   = "email"
	AddressScope = "address"
	PhoneScope   = "phone"
)

type ProfileClaims struct {
	Name              string `json:"name"`
	FamilyName        string `json:"family_name"`
	GivenName         string `json:"given_name"`
	MiddleName        string `json:"middle_name"`
	Nickname          string `json:"nickname"`
	PreferredUsername string `json:"preferred_username"`
	Profile           string `json:"profile"`
	Picture           string `json:"picture"`
	Website           string `json:"website"`
	Gender            string `json:"gender"`
	Birthdate         string `json:"birthdate"`
	ZoneInfo          string `json:"zoneinfo"`
	Locale            string `json:"locale"`
	UpdatedAt         string `json:"updated_at"`
}
type EmailClaims struct {
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}
type AddressClaims struct {
	Address string `json:"address"`
}
type PhoneClaims struct {
	PhoneNumber         string `json:"phone_number"`
	PhoneNumberVerified bool   `json:"phone_number_verified"`
}

type DMF struct {
	Models        []Models `json:"models"`
	Enums         []Enums  `json:"enums"`
	SchemaContent string   `json:"schemaContent"`
}

type Models struct {
	Id            string   `json:"id"`
	Name          string   `json:"name"`
	Create        bool     `json:"create"`
	Delete        bool     `json:"delete"`
	Update        bool     `json:"update"`
	IdField       string   `json:"idField"`
	DisplayFields []string `json:"displayFields"`
	Fields        []Fields `json:"fields"`
}

type Fields struct {
	Id            string      `json:"id"`
	Name          string      `json:"name"`
	Title         string      `json:"title"`
	Type          string      `json:"type"`
	List          bool        `json:"list"`
	Kind          string      `json:"kind"`
	Read          bool        `json:"read"`
	Required      bool        `json:"required"`
	IsId          bool        `json:"isId"`
	Unique        bool        `json:"unique"`
	Create        bool        `json:"create"`
	Order         int         `json:"order"`
	Update        bool        `json:"update"`
	Sort          bool        `json:"sort"`
	Filter        bool        `json:"filter"`
	Editor        bool        `json:"editor"`
	Upload        bool        `json:"upload"`
	RelationField interface{} `json:"relationField"`
}

type Enums struct {
	Name   string   `json:"name"`
	Fields []string `json:"fields"`
}

func GetDMF(dmf dmmf.Document) DMF {
	var dmfModels []Models
	var dmfEnums []Enums
	for _, row := range dmf.Datamodel.Models {
		var models Models
		models.Id = row.Name.String()
		models.Name = row.Name.String()
		models.Create = true
		models.Delete = true
		models.Update = true
		models.IdField = ""
		displayFields := make([]string, 0)
		modelsFields := make([]Fields, 0)
		for i := 0; i < len(row.Fields); i++ {
			fieldRow := row.Fields[i]
			if fieldRow.IsID == true {
				models.IdField = fieldRow.Name.String()
			}
			if fieldRow.Kind == dmmf.FieldKindScalar {
				displayFields = append(displayFields, fieldRow.Name.String())
			}
			if fieldRow.Kind == dmmf.FieldKindObject && len(fieldRow.RelationToFields) == 0 {
				continue
			}
			// Kind如果是object类型则需要遍历该object并将字段添加到displayFields里去
			if fieldRow.Kind == dmmf.FieldKindObject {
				// 数组序列化
				content, _ := json.Marshal(fieldRow.RelationToFields)
				fields := strings.Trim(string(content), "[]")
				fields = strings.ReplaceAll(fields, "\"", "")
				for _, name := range strings.Split(fields, ",") {
					// 此时它到类型就是表名
					displayFields = append(displayFields, fmt.Sprintf("%s.%s", fieldRow.Type, name))
				}
				continue
			}

			//displayFields = append(displayFields, fieldRow.Name.String())
			modelsFields = append(modelsFields, Fields{
				Id:            fmt.Sprintf("%s.%s", row.Name.String(), fieldRow.Name.String()),
				Name:          fieldRow.Name.String(),
				Title:         fieldRow.Name.String(),
				Type:          fieldRow.Type.String(),
				List:          fieldRow.IsList,
				Kind:          string(fieldRow.Kind),
				Read:          true,
				Required:      fieldRow.IsRequired,
				IsId:          fieldRow.IsID,
				Unique:        fieldRow.IsUnique,
				Create:        true,
				Order:         i + 1,
				Update:        true,
				Sort:          true,
				Filter:        true,
				Editor:        true,
				Upload:        true,
				RelationField: fieldRow.RelationToFields,
			})
		}
		models.DisplayFields = displayFields
		models.Fields = modelsFields
		dmfModels = append(dmfModels, models)
	}
	for _, row := range dmf.Datamodel.Enums {
		var enum Enums
		enum.Name = row.Name.String()
		enums := make([]string, 0)
		for _, value := range row.Values {
			enums = append(enums, value.Name.String())
		}
		enum.Fields = enums
		dmfEnums = append(dmfEnums, enum)
	}

	return DMF{
		Models: dmfModels,
		Enums:  dmfEnums,
	}
}
