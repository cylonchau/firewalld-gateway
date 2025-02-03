package model

import (
	"gorm.io/gorm"
)

const port_table_name = "ports"

type Port struct {
	gorm.Model
	Port       uint16 `json:"port" gorm:"index;type:smallint;unsigned"`
	Protocol   string `json:"protocol" gorm:"type:varchar(255)"`
	TemplateId int    `json:"template_id" gorm:"type:int"`
}

type PortList struct {
	ID         int    `json:"id"`
	Port       uint16 `json:"port"`
	Protocol   string `json:"protocol"`
	Template   string `json:"template"`
	TemplateID int    `json:"template_id"`
}

type PortListWithoutID struct {
	Port     string `json:"port"`
	Protocol string `json:"protocol"`
}

func (*PortList) TableName() string {
	return port_table_name
}

func (*PortListWithoutID) TableName() string {
	return port_table_name
}

func CreatePort(port uint16, protocol string, template_id int) (enconterError error) {
	Port := &Port{
		Port:       port,
		Protocol:   protocol,
		TemplateId: template_id,
	}
	result := DB.Create(Port)
	if enconterError = result.Error; enconterError == nil {
		return nil
	}

	return enconterError
}

func GetPorts(title string, offset, limit int, sort string) (map[string]interface{}, error) {
	ports := []*PortList{}
	response := make(map[string]interface{})
	var count int64
	portQuery := DB.Table(port_table_name).
		Select([]string{
			port_table_name + ".id",
			port_table_name + ".port",
			port_table_name + ".protocol",
			"templates.name template",
			"templates.id template_id"}).
		Joins("JOIN templates ON "+port_table_name+".template_id = templates.id").
		Where(port_table_name+".deleted_at IS ?", nil)
	if title != "" {
		portQuery.Where("ports.`port` LIKE ?", title+"%")
	}
	result := portQuery.
		Order(port_table_name + ".id " + sort).
		Limit(limit).Offset((offset - 1) * limit).
		Scan(&ports)
	DB.Model(&Port{}).Distinct("id").Count(&count)
	if result.Error != gorm.ErrRecordNotFound {
		response["list"] = ports
		response["total"] = count
		return response, nil
	}
	return nil, result.Error
}

func CheckPortIsExistWithName(name string) bool {
	result := DB.Where("name = ?", name).First(&Port{})
	if result.Error != gorm.ErrRecordNotFound || result.RowsAffected > 0 {
		return false
	}
	return true
}

func DeletePortWithID(id uint64) error {
	result := DB.Delete(&Port{}, id)
	if result.Error == nil {
		return nil
	}
	return result.Error
}

func UpdatePortWithID(id uint64, port uint16, protocol string, template_id int) (enconterError error) {
	p := &Port{
		Port:       port,
		Protocol:   protocol,
		TemplateId: template_id,
	}
	result := DB.Model(&Port{}).Where("id = ?", id).Updates(p)
	if enconterError = result.Error; enconterError == nil {
		return nil
	}

	return enconterError
}
