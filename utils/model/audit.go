package model

import (
	"gorm.io/gorm"
)

const audit_table_name = "audits"

type Audit struct {
	gorm.Model
	UserID  uint64 `json:"user_id" gorm:"index;type:int"`
	IP      uint32 `json:"ip" gorm:"index;type:int"`
	Method  string `json:"method" gorm:"type:char(5)"`
	Path    string `json:"path" gorm:"varchar(50)"`
	Browser string `json:"browser" gorm:"varchar(50)"`
	System  string `json:"system" gorm:"varchar(50)"`
}

type AuditList struct {
	Username string `json:"username" gorm:"index;type:varchar(20)"`
	IP       uint32 `json:"ip" gorm:"index;type:int"`
	Method   string `json:"method" gorm:"type:char(5)"`
	Path     string `json:"path" gorm:"varchar(50)"`
	Browser  string `json:"browser" gorm:"varchar(50)"`
	System   string `json:"system" gorm:"varchar(50)"`
}

func (*Audit) TableName() string {
	return audit_table_name
}

func (*AuditList) TableName() string {
	return audit_table_name
}

func GetAuditLogs(title string, offset, limit int, sort string) (map[string]interface{}, error) {
	logs := []*AuditList{}
	response := make(map[string]interface{})
	var count int64

	// 查询审计日志
	query := DB.Table(audit_table_name).
		Select([]string{
			audit_table_name + ".ip",
			audit_table_name + ".path",
			audit_table_name + ".method",
			audit_table_name + ".browser",
			audit_table_name + ".system",
			user_table_name + ".username",
		}).
		Joins("join "+user_table_name+" on "+user_table_name+".id = "+audit_table_name+".user_id").
		Where(audit_table_name+".deleted_at IS ?", nil).
		Limit(limit).
		Offset((offset - 1) * limit)

	if title != "" {
		query = query.Where(user_table_name+".username LIKE ?", "%"+title+"%")
	}

	result := query.Order(audit_table_name + ".id " + sort).Scan(&logs)

	// 查询总数
	totalQuery := DB.Table(audit_table_name).
		Where(audit_table_name+".deleted_at IS ?", nil)

	if title != "" {
		totalQuery = totalQuery.Joins("join "+user_table_name+" on "+user_table_name+".id = "+audit_table_name+".user_id").
			Where(user_table_name+".username LIKE ?", "%"+title+"%")
	}

	// 获取内容的总数
	totalQuery.Count(&count)

	if result.Error != gorm.ErrRecordNotFound {
		response["list"] = logs
		response["total"] = count
		return response, nil
	}
	return nil, result.Error
}

func AppendAuditLog(auditLog map[string]interface{}) {

	auditItem := &Audit{
		UserID:  uint64(auditLog["user_id"].(int64)),
		IP:      auditLog["ip"].(uint32), // 需要类型转换
		Method:  auditLog["method"].(string),
		Path:    auditLog["path"].(string),
		Browser: auditLog["browser"].(string),
		System:  auditLog["system"].(string),
	}
	DB.Create(auditItem)
}
