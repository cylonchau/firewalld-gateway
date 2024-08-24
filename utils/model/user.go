package model

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/praserx/ipconv"
	"gorm.io/gorm"

	"github.com/cylonchau/firewalld-gateway/config"
	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
)

const Secret = "com.github.cylonchau"

var user_table_name = "users"

type User struct {
	gorm.Model
	Username string `gorm:"index;type:varchar(20)"`
	Password string `gorm:"type:varchar(32)"`
	Roles    []Role `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	LoginIP  int    `json:"login_ip" gorm:"index;type:int"`
}

type UserInfo struct {
	ID       uint       `gorm:"primarykey;column:id" json:"id"`
	Username string     `gorm:"index;type:varchar(20);column:username" json:"username"`
	Roles    []RoleInfo `gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:ID;joinForeignKey:UserID;references:ID;joinReferences:RoleID;column:roles" json:"roles"`
}

type UserList struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"update_at"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index"`
	Username  string         `json:"username" gorm:"index;type:varchar(20)"`
	Roles     []Role         `json:"roles" gorm:"many2many:user_roles;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	LoginIP   int            `json:"login_ip" gorm:"index;type:int"`
}

func (*User) TableName() string {
	return user_table_name
}

func (*UserInfo) TableName() string {
	return user_table_name
}

func (*UserList) TableName() string {
	return user_table_name
}

func QueryUserWithUsername(username string) (User, error) {
	user := User{}
	result := DB.Select("id", "username", "password").Where("username = ?", username).Find(&user)
	if result.Error == nil {
		return user, nil
	}
	return User{}, result.Error
}

func QueryUserWithUID(uid int64) (User, error) {
	user := User{}
	result := DB.Select("id", "username").Where("id = ?", uid).Find(&user)
	if result.Error == nil {
		return user, nil
	}
	return User{}, result.Error
}

func CreateUser(userQuery *query.UserQuery) (enconterError error) {
	if checkUserExist(userQuery.Username) {
		user := &User{
			Username: userQuery.Username,
			Password: EncryptPassword(userQuery.Password),
		}
		result := DB.Create(user)
		if enconterError = result.Error; enconterError == nil {
			return nil
		}
	} else {
		enconterError = query.ErrUserExist
	}
	return enconterError
}

func UpdateUserWithID(query *query.UserEditQuery) (enconterError error) {
	user := &User{
		Username: query.Username,
		Password: EncryptPassword(query.Password),
	}

	result := DB.Model(&User{Model: gorm.Model{ID: uint(query.ID)}}).Updates(user)
	if enconterError = result.Error; enconterError == nil {
		return nil
	}

	return enconterError
}

func AllocateRole2User(userID uint64, roles_id []int) (enconterError error) {
	roles := GenerateRoleWithID(roles_id)
	if enconterError = DB.Model(&User{Model: gorm.Model{ID: uint(userID)}}).Association("Roles").Replace(roles); enconterError == nil || enconterError != gorm.ErrRecordNotFound {
		return nil
	}
	return enconterError
}

func LastLogin(uid int64, ip uint32) bool {
	result := DB.Model(&User{}).Where("id = ?", uid).Update("login_ip", ip)
	if result.Error == nil {
		return true
	}
	return false
}

func EncryptPassword(p string) string {
	h := md5.New()
	h.Write([]byte(Secret + config.CONFIG.AppName))
	h.Sum([]byte(p))
	return hex.EncodeToString(h.Sum([]byte(p)))
}

func checkUserExist(username string) bool {
	result := DB.Where("username = ?", username).First(&User{})
	if result.Error != gorm.ErrRecordNotFound || result.RowsAffected > 0 {
		return false
	}
	return true
}

func GetRequestIP(r *http.Request) (uint32, error) {
	ip := r.Header.Get("X-Real-IP")
	if net.ParseIP(ip) != nil {
		return 0, nil
	}

	ip = r.Header.Get("X-Forward-For")
	for _, i := range strings.Split(ip, ",") {
		if net.ParseIP(i) != nil {
			return 0, nil
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return 0, err
	}

	if ip, version, err := ipconv.ParseIP(ip); err == nil && version == 4 {
		return ipconv.IPv4ToInt(ip)
	}

	return 0, errors.New("no valid ip found")
}

func GetUsers(offset, limit int, sort string) (map[string]interface{}, error) {
	users := []*UserList{}
	response := make(map[string]interface{})
	var count int64
	result := DB.
		Limit(limit).
		Offset(offset).
		Where("deleted_at is ?", nil).
		Order("id " + sort).
		Find(&users)
	DB.Model(&User{}).Distinct("id").Count(&count)
	if result.Error != gorm.ErrRecordNotFound {
		response["list"] = users
		response["total"] = count
		return response, nil
	}
	return nil, result.Error
}

func DeleteUserWithID(id uint64) error {
	result := DB.Delete(&User{}, id)
	if result.Error == nil && result.RowsAffected > 0 {
		return nil
	}
	return result.Error
}
