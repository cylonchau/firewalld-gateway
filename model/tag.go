package model

import (
	"fmt"

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
		enconterError = apis.ErrUserTagExist
	}
	return enconterError
}

func GetTags(offset, limit int) ([]TagInfo, error) {
	tags := []TagInfo{}
	result := DB.Select([]string{"id", "name"}).Limit(limit).Offset(offset).Where("deleted_at is ?", nil).Find(&tags)
	if result.Error != gorm.ErrRecordNotFound {
		return tags, nil
	}
	return nil, result.Error
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

	fmt.Printf("%+v", result)
	if result.Error == nil && result.RowsAffected > 0 {
		return nil
	}
	return result.Error
}
