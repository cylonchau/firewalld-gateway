package model

import (
	"gorm.io/gorm"

	"github.com/cylonchau/firewalld-gateway/server/apis"
)

type Tag struct {
	gorm.Model
	Name        string `json:"name" gorm:"index;type:varchar(255)"`
	Description string `json:"description" gorm:"type:varchar(255)"`
	Hosts       []Host `json:"hosts" gorm:"foreignKey:TagId;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type TagInfo struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (*Tag) TableName() string {
	return "tags"
}

func (*TagInfo) TableName() string {
	return "tags"
}

func CreateTag(tagQuery *apis.TagEditQuery) (enconterError error) {
	if CheckTagIsExistWithName(tagQuery.Name) {
		tag := &Tag{
			Name:        tagQuery.Name,
			Description: tagQuery.Description,
		}
		result := DB.Create(tag)
		if enconterError = result.Error; enconterError == nil {
			return nil
		}
	} else {
		enconterError = apis.ErrTagExist
	}
	return enconterError
}

func GetTags(offset, limit int, sort string) (map[string]interface{}, error) {
	tags := []*TagInfo{}
	response := make(map[string]interface{})
	var count int64
	result := DB.Select([]string{"id", "name", "description"}).
		Limit(limit).
		Offset(offset).
		Where("deleted_at is ?", nil).
		Order("tags.id " + sort).
		Find(&tags)
	DB.Model(&Tag{}).Distinct("id").Count(&count)
	if result.Error != gorm.ErrRecordNotFound {
		response["list"] = tags
		response["total"] = count
		return response, nil
	}
	return nil, result.Error
}

func UpdateTagWithID(query *apis.TagEditQuery) (enconterError error) {
	tag := &Tag{
		Name:        query.Name,
		Description: query.Description,
	}
	result := DB.Model(&Tag{}).Where("id = ?", query.ID).Updates(tag)
	if enconterError = result.Error; enconterError == nil {
		return nil
	}

	return enconterError
}

func CheckTagIsExistWithID(id int64) (Tag, error) {
	tag := Tag{}
	result := DB.Select("id", "name").Where("id = ?", id).Find(&tag)
	if result.Error == nil {
		return tag, nil
	}
	return Tag{}, result.Error
}

func CheckTagIsExistWithName(name string) bool {
	result := DB.Where("name = ?", name).First(&Tag{})
	if result.Error != gorm.ErrRecordNotFound || result.RowsAffected > 0 {
		return false
	}
	return true
}

func DeleteTagWithID(id uint64) error {
	result := DB.Delete(&Tag{}, id)
	if result.Error == nil && result.RowsAffected > 0 {
		return nil
	}
	return result.Error
}
