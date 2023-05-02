package model

import (
	"net"

	"github.com/praserx/ipconv"
	"gorm.io/gorm"

	"github.com/cylonchau/firewalld-gateway/server/apis"
)

var host_table_name = "hosts"

type Host struct {
	gorm.Model
	Hostname string `json:"hostname" gorm:"index;type:varchar(255)"`
	IP       uint32 `json:"ip" gorm:"index;type:int"`
	TagId    int    `json:"tag_id" gorm:"index;type:int"`
}

type HostList struct {
	ID       int    `json:"id"`
	Hostname string `json:"hostname"`
	Ip       uint32 `json:"ip"`
	Tag      string `json:"tag"`
	TagId    int    `json:"tag_id"`
}

func (*HostList) TableName() string {
	return host_table_name
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

func UpdateHostWithID(query *apis.HostQuery) (enconterError error) {
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
			IP:       ip,
			Hostname: query.Hostname,
			TagId:    query.TagId,
		}
		result := DB.Model(&Host{}).Where("id = ?", query.ID).Updates(host)
		if enconterError = result.Error; enconterError == nil {
			return nil
		}
	}
	return enconterError
}

func CreateHostWithHost(host *Host) (enconterError error) {
	result := DB.Create(host)
	if enconterError = result.Error; enconterError == nil {
		return nil
	}

	return enconterError
}

func CreateHost(hostIP, host string, tagID int) (enconterError error) {
	var (
		ip      uint32
		version int
		netip   net.IP
	)
	if netip, version, enconterError = ipconv.ParseIP(hostIP); enconterError == nil && version == 4 {
		if ip, enconterError = ipconv.IPv4ToInt(netip); enconterError != nil {
			return enconterError
		}
		host := &Host{
			IP:       ip,
			Hostname: host,
			TagId:    tagID,
		}
		result := DB.Create(host)
		if enconterError = result.Error; enconterError == nil {
			return nil
		}
	}
	return enconterError
}

func GetHosts(offset, limit int, sort string) (map[string]interface{}, error) {
	hosts := []*HostList{}
	var count int64
	response := make(map[string]interface{})
	result := DB.Table(host_table_name).
		Select([]string{
			host_table_name + ".id",
			host_table_name + ".hostname",
			host_table_name + ".ip",
			"tags.name tag",
			"tags.id tag_id"}).
		Joins("join tags on "+host_table_name+".tag_id = tags.id").
		Limit(limit).
		Offset((offset-1)*limit).
		Where(host_table_name+".deleted_at is ?", nil).
		Order(host_table_name + ".id " + sort).
		Scan(&hosts)
	DB.Model(&Host{}).Distinct("id").Count(&count)
	if result.Error != gorm.ErrRecordNotFound {
		response["list"] = hosts
		response["total"] = count
		return response, nil
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
