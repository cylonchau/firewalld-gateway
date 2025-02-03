package model

import (
	"net"

	"github.com/praserx/ipconv"
	"gorm.io/gorm"

	"github.com/cylonchau/firewalld-gateway/utils/apis/query"
)

const host_table_name = "hosts"

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

type Classify struct {
	Name  string `json:"name"`
	Count int    `json:"value"`
}

func (*HostList) TableName() string {
	return host_table_name
}

func (*Classify) TableName() string {
	return host_table_name
}

func GetHostsByTagName(tagName string) ([]Host, error) {
	var hosts []Host
	// 查询指定名称的 Tag，并预加载 Hosts
	if err := DB.Table(host_table_name).Joins("JOIN tags ON tags.id = hosts.tag_id").
		Where("tags.name = ?", tagName).
		Find(&hosts).Error; err != nil {
		return nil, err
	}

	// 返回与该 Tag 相关的 Hosts
	return hosts, nil
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

func UpdateHostWithID(query *query.HostQuery) (enconterError error) {
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

func HostCounter() int64 {
	var count int64
	DB.Model(&Host{}).Distinct("id").Count(&count)
	return count
}

func HostClassify() ([]*Classify, error) {
	hosts := []*Classify{}
	result := DB.Table(host_table_name).
		Joins("join tags on " + host_table_name + ".tag_id = tags.id").
		Select([]string{
			"COUNT(*) AS count",
			"tags.name AS name"}).
		Group("tags.name").
		Where(host_table_name + ".`deleted_at` IS NULL").
		Where("tags.`deleted_at` IS NULL").
		Scan(&hosts)
	if result.Error != gorm.ErrRecordNotFound || result.Error == nil {
		return hosts, nil
	}
	return nil, result.Error
}
