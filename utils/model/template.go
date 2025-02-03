package model

import (
	"github.com/jinzhu/copier"
	"gorm.io/gorm"

	"github.com/cylonchau/firewalld-gateway/api"
	"github.com/cylonchau/firewalld-gateway/config"
	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
)

const template_table_name = "templates"

type Template struct {
	gorm.Model
	Name        string `json:"name" gorm:"index;type:varchar(255)"`
	Description string `json:"description" gorm:"type:varchar(255)"`
	Target      string `json:"target" gorm:"type:varchar(100)"`
}

type TemplateList struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Target      string `json:"target,omitempty"`
	Description string `json:"description,omitempty"`
}

type TemplateListWithoutID struct {
	Name        string `json:"name"`
	Target      string `json:"target,omitempty"`
	Description string `json:"description,omitempty"`
}

type TemplateWithDetails struct {
	Target      string              `json:"target"`
	Description string              `json:"description"`
	Short       string              `json:"short"`
	Riches      []RichListWithoutID `json:"rich"`
	Ports       []PortListWithoutID `json:"port"`
}

func (*TemplateList) TableName() string {
	return template_table_name
}

func (*TemplateListWithoutID) TableName() string {
	return template_table_name
}

func GetRichWithDetailsByTemplateID(templateID uint) (*api.Settings, error) {
	result := &api.Settings{}

	// 查询 Template
	template := &Template{}
	if err := DB.Table(template_table_name).First(&template, templateID).Error; err != nil {
		return nil, err
	}
	result.Target = template.Target
	result.Description = template.Description
	result.Short = template.Name
	// 查询所有 Rich 记录
	var riches []RichListWithoutID
	if err := DB.Where("template_id = ?", templateID).Where(rich_table_name+".deleted_at is ?", nil).Find(&riches).Error; err != nil {
		return nil, err
	}

	apiRichRule := []*api.Rule{}
	copier.Copy(&apiRichRule, &riches)
	for k, v := range apiRichRule {
		switch riches[k].Action {
		case "accept":
			v.Accept = &api.Accept{Flag: true}
		case "drop":
			v.Drop = &api.Drop{Flag: true}
		case "reject":
			v.Reject = &api.Reject{Flag: true}
		}
		result.Rule = append(result.Rule, v.ToString())
	}
	//
	//// 查询所有 Port 记录
	var ports []*PortListWithoutID
	if err := DB.Where("template_id = ?", templateID).Where(port_table_name+".deleted_at is ?", nil).Find(&ports).Error; err != nil {
		return result, err
	}

	apiPortsRule := []*api.Port{}
	copier.Copy(&apiPortsRule, &ports)
	result.Port = apiPortsRule

	result.Port = append(result.Port, &api.Port{
		Port:     config.CONFIG.DbusPort,
		Protocol: "tcp",
	})

	return result, nil
}

func CreateTemplate(name, description, target string) (enconterError error) {
	if CheckTemplateIsExistWithName(name) {
		template := &Template{
			Name:        name,
			Description: description,
			Target:      target,
		}
		result := DB.Create(template)
		if enconterError = result.Error; enconterError == nil {
			return nil
		}
	} else {
		enconterError = query.ErrTemplateExist
	}
	return enconterError
}

func GetTemplates(title string, offset, limit int, sort string) (map[string]interface{}, error) {
	templates := []*TemplateList{}
	response := make(map[string]interface{})
	var count int64
	result := DB.Select([]string{"id", "name", "description", "target"}).
		Limit(limit).Offset(offset).
		Where("deleted_at is ?", nil).
		Where(template_table_name+".name LIKE ?", "%"+title+"%").
		Order(template_table_name + ".id " + sort).
		Find(&templates)
	DB.Model(&Template{}).Distinct("id").Count(&count)
	if result.Error != gorm.ErrRecordNotFound {
		response["list"] = templates
		response["total"] = count
		return response, nil
	}
	return nil, result.Error
}

func GetSimpleTemplates(offset, limit int, sort string) (map[string]interface{}, error) {
	templates := []*TemplateList{}
	response := make(map[string]interface{})
	result := DB.Select([]string{"id", "name"}).
		Limit(limit).Offset(offset).
		Where("deleted_at is ?", nil).
		Order(template_table_name + ".id " + sort).
		Find(&templates)
	if result.Error != gorm.ErrRecordNotFound {
		response["list"] = templates
		return response, nil
	}
	return nil, result.Error
}

func CheckTemplateIsExistWithName(name string) bool {
	result := DB.Where("name = ?", name).First(&Template{})
	if result.Error != gorm.ErrRecordNotFound || result.RowsAffected > 0 {
		return false
	}
	return true
}

func DeleteTemplateWithID(id uint64) error {
	result := DB.Delete(&Template{}, id)
	if result.Error == nil {
		return nil
	}
	return result.Error
}

func UpdateTemplateWithID(id uint64, name, description, target string) (enconterError error) {
	template := &Template{
		Name:        name,
		Description: description,
		Target:      target,
	}
	result := DB.Model(&Template{}).Where("id = ?", id).Updates(template)
	if enconterError = result.Error; enconterError == nil {
		return nil
	}

	return enconterError
}

func TemplateCounter() int64 {
	var count int64
	DB.Model(&Template{}).Distinct("id").Count(&count)
	return count
}
