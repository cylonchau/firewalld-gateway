package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"strings"

	json "github.com/json-iterator/go"
	"gorm.io/gorm"

	"github.com/cylonchau/firewalld-gateway/api"
	queryapi "github.com/cylonchau/firewalld-gateway/utils/apis/query"
)

const rich_table_name = "riches"

type Rich struct {
	gorm.Model
	Family      string           `form:"family" json:"family,omitempty" gorm:"type:char(4)"`
	Source      *api.Source      `form:"source" json:"source,omitempty" gorm:"json"`
	Destination *api.Destination `form:"destination" json:"destination,omitempty" gorm:"json"`
	Port        *api.Port        `form:"port" json:"port,omitempty" gorm:"json"`
	Protocol    *api.Protocol    `form:"protocol" json:"protocol,omitempty" gorm:"json"`
	Action      string           `form:"action" json:"action,omitempty" gorm:"type:varchar(30)"`
	Limit       uint16           `form:"limit" json:"limit,omitempty" gorm:"type:varchar(255)"`
	LimitUnit   string           `form:"limit_unit" json:"limit_unit,omitempty" gorm:"type:char(1)"`
	TemplateID  int              `form:"template_id" json:"template_id,omitempty" gorm:"type:int"`
}

type RichList struct {
	ID          int              `json:"id"`
	Family      string           `form:"family" json:"family,omitempty"`
	Source      *api.Source      `form:"source" json:"source,omitempty" gorm:"json"`
	Destination *api.Destination `form:"destination" json:"destination,omitempty" gorm:"json"`
	Port        *api.Port        `form:"port" json:"port,omitempty" gorm:"json"`
	Protocol    *api.Protocol    `form:"protocol" json:"protocol,omitempty" gorm:"json"`
	Action      string           `form:"action" json:"action,omitempty"`
	Limit       uint16           `form:"limit" json:"limit,omitempty"`
	LimitUnit   string           `form:"limit_unit" json:"limit_unit,omitempty"`
	TemplateID  int              `form:"template_id" json:"template_id,omitempty"`
	Template    string           `json:"template"`
}

type RichListWithoutID struct {
	Family      string           `form:"family" json:"family,omitempty"`
	Source      *api.Source      `form:"source" json:"source,omitempty" gorm:"json"`
	Destination *api.Destination `form:"destination" json:"destination,omitempty" gorm:"json"`
	Port        *api.Port        `form:"port" json:"port,omitempty" gorm:"json"`
	Protocol    *api.Protocol    `form:"protocol" json:"protocol,omitempty" gorm:"json"`
	Action      string           `form:"action" json:"action,omitempty"`
	Limit       uint16           `form:"limit" json:"limit,omitempty"`
	LimitUnit   string           `form:"limit_unit" json:"limit_unit,omitempty"`
}

func (r *Rich) Scan(value interface{}) error {
	b, _ := value.([]byte)
	return json.Unmarshal(b, r)
}

func (*RichList) TableName() string {
	return rich_table_name
}

func (*RichListWithoutID) TableName() string {
	return rich_table_name
}

func (r *Rich) Value() (value driver.Value, err error) {
	return json.Marshal(r)
}

func (r *RichList) Scan(value interface{}) error {
	var encounterError error
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New(fmt.Sprint("Failed to unmarshal JSONB value:", value))
	}

	if len(bytes) > 2 {
		var result = &RichList{}
		encounterError = json.Unmarshal(bytes, result)
		if encounterError == nil {
			panic(result)
			if !isEmptyStruct(result.Port) {

				r.Port = result.Port
			}
			if !(result.Source == &api.Source{}) {
				r.Source = result.Source
			}

			if !(result.Destination == &api.Destination{}) {
				r.Destination = result.Destination
			}

			if (result.Protocol == &api.Protocol{}) {
				r.Protocol = result.Protocol
			}
		}
	}
	return encounterError
}

func (r *RichList) Value() (value driver.Value, err error) {
	return json.Marshal(r)
}

func CreateRich(query *queryapi.RichEditQuery) (enconterError error) {

	var port *api.Port
	if query.Port != "" {
		s := strings.Split(query.Port, "/")
		port = &api.Port{Port: s[0], Protocol: s[1]}
	}

	rich := &Rich{
		Family: query.Family,
		Source: &api.Source{
			Address: query.Source,
		},
		Destination: &api.Destination{
			Address: query.Destination,
		},
		Port: port,
		Protocol: &api.Protocol{
			Value: query.Protocol,
		},
		Action:     query.Action,
		Limit:      query.Limit,
		LimitUnit:  query.LimitUnit,
		TemplateID: query.TemplateId,
	}
	result := DB.Create(rich)
	if enconterError = result.Error; enconterError == nil {
		return nil
	}

	return enconterError
}

func GetRich(title string, offset, limit int, sort string) (map[string]interface{}, error) {
	richs := []*RichList{}
	response := make(map[string]interface{})
	var count int64
	richQuery := DB.Table(rich_table_name).
		Select([]string{
			rich_table_name + ".id",
			rich_table_name + ".family",
			rich_table_name + ".source",
			rich_table_name + ".destination",
			rich_table_name + ".port",
			rich_table_name + ".protocol",
			rich_table_name + ".action",
			rich_table_name + ".`limit`",
			rich_table_name + ".limit_unit",
			"templates.name template",
			"templates.id template_id"}).
		Joins("JOIN templates ON "+rich_table_name+".template_id = templates.id").
		Where(rich_table_name+".deleted_at is ?", nil)
	if title != "" {
		richQuery.Where(rich_table_name+".`destination` LIKE ?", "%"+title+"%")
	}
	result := richQuery.Order(rich_table_name + ".id " + sort).
		Limit(limit).
		Offset((offset - 1) * limit).
		Scan(&richs)
	DB.Model(&Rich{}).Distinct("id").Count(&count)
	if result.Error != gorm.ErrRecordNotFound {
		response["list"] = richs
		response["total"] = count
		return response, nil
	}
	return nil, result.Error
}

func isEmptyStruct(s interface{}) bool {
	v := reflect.ValueOf(s)
	for i := 0; i < v.NumField(); i++ {
		fv := v.Field(i)
		switch fv.Kind() {
		case reflect.String:
			if fv.String() != "" {
				return false
			}
		case reflect.Slice, reflect.Map:
			if !fv.IsNil() && fv.Len() > 0 {
				return false
			}
		case reflect.Bool:
			if fv.Bool() {
				return false
			}
		default:
			if !fv.IsZero() {
				return false
			}
		}
	}
	return true
}

func UpdateRichWithID(query *queryapi.RichEditQuery) (enconterError error) {
	var port *api.Port
	if query.Port != "" {
		s := strings.Split(query.Port, "/")
		port = &api.Port{Port: s[0], Protocol: s[1]}
	}

	rich := &Rich{
		Family: query.Family,
		Source: &api.Source{
			Address: query.Source,
		},
		Destination: &api.Destination{
			Address: query.Destination,
		},
		Port: port,
		Protocol: &api.Protocol{
			Value: query.Protocol,
		},
		Action:     query.Action,
		Limit:      query.Limit,
		LimitUnit:  query.LimitUnit,
		TemplateID: query.TemplateId,
	}
	result := DB.Model(&Rich{}).Where("id = ?", query.ID).Updates(rich)
	if enconterError = result.Error; enconterError == nil {
		return nil
	}

	return enconterError
}

func DeleteRichWithID(id uint64) error {
	result := DB.Delete(&Rich{}, id)
	if result.Error == nil {
		return nil
	}
	return result.Error
}
