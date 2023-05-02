package model

import (
	"gorm.io/gorm"

	"github.com/cylonchau/firewalld-gateway/server/apis"
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

func (*TemplateList) TableName() string {
	return template_table_name
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
		enconterError = apis.ErrTemplateExist
	}
	return enconterError
}

func GetTemplates(offset, limit int, sort string) (map[string]interface{}, error) {
	templates := []*TemplateList{}
	response := make(map[string]interface{})
	var count int64
	result := DB.Select([]string{"id", "name", "description", "target"}).
		Limit(limit).Offset(offset).
		Where("deleted_at is ?", nil).
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
