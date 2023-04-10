package model

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/cylonchau/firewalld-gateway/config"
	"github.com/cylonchau/firewalld-gateway/server/apis"
	"github.com/praserx/ipconv"
	"gorm.io/gorm"
)

const Secret = "com.github.cylonchau"

type User struct {
	gorm.Model
	Username string `gorm:"index;type:varchar(20)"`
	Password string `gorm:"type:varchar(32)"`
	LoginIP  int    `json:"login_ip" gorm:"type:int"`
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

func CreateUser(userQuery *apis.UserQuery) (enconterError error) {
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
		enconterError = apis.ErrUserExist
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

	if ip, version, err := ipconv.ParseIP(ip); err != nil && version == 4 {
		return ipconv.IPv4ToInt(ip)
	}

	return 0, errors.New("no valid ip found")
}
