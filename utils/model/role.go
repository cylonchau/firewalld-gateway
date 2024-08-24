package model

import (
	"gorm.io/gorm"

	query2 "github.com/cylonchau/firewalld-gateway/utils/apis/query"
)

var role_table_name = "roles"

// Role model
type Role struct {
	gorm.Model
	Name        string   `gorm:"size:50;not null;unique" json:"name"`
	Routers     []Router `gorm:"many2many:role_routers;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Users       []User   `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Tokens      []Token  `gorm:"many2many:token_roles;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Description string   `gorm:"size:255;not null" json:"description"`
}

type RoleList struct {
	ID          int    `json:"id"`
	Name        string `gorm:"size:50;not null;unique" json:"name"`
	Description string `gorm:"size:255;not null" json:"description"`
}

type RoleInfo struct {
	ID   int    `json:"id" gorm:"column:id"`
	Name string `json:"name" gorm:"column:name"`
}

type UserRoles struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func (*Role) TableName() string {
	return role_table_name
}

func (*RoleInfo) TableName() string {
	return role_table_name
}

func (*RoleList) TableName() string {
	return role_table_name
}

func (*UserRoles) TableName() string {
	return role_table_name
}

// Get all roles
func GetRoles(offset, limit int, sort string) (map[string]interface{}, error) {
	roles := []*RoleList{}
	response := make(map[string]interface{})
	var count int64
	result := DB.Select([]string{"id", "name", "description"}).
		Limit(limit).
		Offset(offset).
		Order("id "+sort).
		Where(role_table_name+".deleted_at is ?", nil).
		Find(&roles)
	DB.Model(&Role{}).Distinct("id").Count(&count)
	if result.Error != gorm.ErrRecordNotFound {
		response["list"] = roles
		response["total"] = count
		return response, nil
	}
	return nil, result.Error
}

// Create role
func CreateRole(query *query2.RoleEditQuery, routers []Router) (enconterError error) {
	if checkRoleIsExistWithName(query.Name) {
		role := &Role{
			Name:        query.Name,
			Description: query.Description,
			Routers:     routers,
		}
		result := DB.Create(role)
		if enconterError := result.Error; enconterError == nil {
			return nil
		}
	} else {
		enconterError = query2.ErrRoleExist
	}
	return enconterError
}

func checkRoleIsExistWithName(name string) bool {
	result := DB.Where("name = ?", name).First(&Role{})
	if result.Error != gorm.ErrRecordNotFound || result.RowsAffected > 0 {
		return false
	}
	return true
}

func UpdateRoleWithID(query *query2.RoleEditQuery) (enconterError error) {
	role := &Role{
		Name:        query.Name,
		Description: query.Description,
	}
	routers := GenerateRouterWithID(query.RouterIDs)
	role.Routers = routers
	result := DB.Model(&Role{}).Where("id = ?", query.ID).Updates(role)
	if enconterError = result.Error; enconterError == nil {
		return nil
	}

	return enconterError
}

func DeleteRoleWithID(id uint64) error {
	result := DB.Delete(&Role{}, id)
	if result.Error == nil && result.RowsAffected > 0 {
		return nil
	}
	return result.Error
}

func GenerateRoleWithID(ids []int) []Role {
	var roleSlice []Role
	DB.Find(&roleSlice, ids)
	return roleSlice
}

func GetRolesWithUID(uid int) (user *UserInfo, encounterError error) {
	DB.Preload("Roles").First(&user, uid)
	encounterError = DB.Model(&User{Model: gorm.Model{ID: uint(uid)}}).
		//Where(role_table_name+".deleted_at is ?", nil).
		Association("Roles").
		Find(&user)
	if encounterError == gorm.ErrRecordNotFound || encounterError == nil {
		return user, nil
	}
	return nil, encounterError
}
