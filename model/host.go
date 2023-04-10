package model

import (
	"errors"
	"net"

	"github.com/praserx/ipconv"
	"gorm.io/gorm"

	"github.com/cylonchau/firewalld-gateway/server/apis"
)

type Host struct {
	gorm.Model
	Hostname string `json:"hostname" gorm:"index;type:varchar(255)"`
	Ip       uint32 `json:"ip" gorm:"index;type:int"`
	TagId    int    `json:"tag_id" gorm:"index;type:int"`
}

type HostList struct {
	ID       int    `json:"id"`
	Hostname string `json:"hostname"`
	Ip       uint32 `json:"ip"`
	TagId    int    `json:"tag_id"`
}

func (*HostList) TableName() string {
	return "hosts"
}

func QueryHostWithName(hostname string) (*Host, error) {
	host := &Host{}
	result := DB.Select("id", "hostname", "ip").Where("hostname = ?", hostname).Find(host)
	if result.Error == nil {
		return host, nil
	}
	return nil, result.Error
}

func QueryHostWithID(uid int) (*Host, error) {
	host := &Host{}
	result := DB.Select("id", "hostname", "ip").Where("id = ?", uid).Find(host)
	if result.Error == nil {
		return host, nil
	}
	return nil, result.Error
}

func CreateHostWithHost(host *Host) (enconterError error) {
	result := DB.Create(host)
	if enconterError = result.Error; enconterError == nil {
		return nil
	}

	return enconterError
}

func CreateHost(query *apis.HostQuery) (enconterError error) {
	var (
		ip      uint32
		version int
		netip   net.IP
	)
	if netip, version, enconterError = ipconv.ParseIP(query.IP); enconterError == nil && version == 4 {
		if ip, enconterError = ipconv.IPv4ToInt(netip); enconterError != nil {
			return enconterError
		}
		host := &Host{
			Ip: ip,
		}
		if query.TagId != 0 {
			host.TagId = query.TagId
		} else {
			return errors.New("tag_id invaild")
		}
		result := DB.Create(host)
		if enconterError = result.Error; enconterError == nil {
			return nil
		}
	}
	return enconterError
}

func GetHosts(offset, limit int) ([]HostList, error) {
	hosts := []HostList{}
	result := DB.Select([]string{"id", "hostname", "ip", "tag_id"}).Limit(limit).Offset(offset).Where("deleted_at is ?", nil).Find(&hosts)
	if result.Error != gorm.ErrRecordNotFound {
		return hosts, nil
	}
	return nil, result.Error
}

func DeleteHostWithID(id uint64) error {
	result := DB.Delete(&Host{}, id)
	if result.Error == nil {
		return nil
	}
	return result.Error
}
