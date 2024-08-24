package model

import (
	"time"

	"gorm.io/gorm"
)

var router_table_name = "routers"

// Router model
type Router struct {
	gorm.Model
	Path   string `json:"path" gorm:"type:varchar(255);not null;"`
	Method string `json:"method" gorm:"type:char(4);not null;"`
	Roles  []Role `gorm:"many2many:role_routers;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type RouterList struct {
	ID        int       `json:"id" gorm:"primarykey;column:id"`
	Path      string    `json:"path" gorm:"type:varchar(255);not null;column:path"`
	Method    string    `json:"method" gorm:"type:char(4);not null;column:method"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at"`
}

func (*Router) TableName() string {
	return router_table_name
}

func (*RouterList) TableName() string {
	return router_table_name
}

// Get all roles
func GetRouters(offset, limit int, sort string) (map[string]interface{}, error) {
	roles := []*RouterList{}
	response := make(map[string]interface{})
	var count int64
	result := DB.Select([]string{"id", "path", "method"}).
		Limit(limit).
		Offset((offset - 1) * limit).
		Order("id " + sort).
		Find(&roles)
	DB.Model(&Router{}).Distinct("id").Count(&count)
	if result.Error != gorm.ErrRecordNotFound {
		response["list"] = roles
		response["total"] = count
		return response, nil
	}
	return nil, result.Error
}

func GenerateRouterWithID(ids []int) []Router {
	var routerSlice []Router
	DB.Find(&routerSlice, ids)
	return routerSlice
}

func GetRoutersWithRID(roleIDs []string) (routers []*RouterList, encounterError error) {
	result := DB.Table("role_routers").
		Joins("JOIN routers ON routers.id = role_routers.router_id").
		Where("role_routers.role_id IN ?", roleIDs).
		Select("routers.id, routers.path, routers.method").
		Scan(&routers)
	if result.Error != gorm.ErrRecordNotFound || result.Error == nil {

		return routers, nil
	}
	return nil, encounterError
}
