package model

import (
	"gorm.io/gorm"

	"github.com/cylonchau/firewalld-gateway/config"
	token2 "github.com/cylonchau/firewalld-gateway/utils/auther"
)

const token_table_name = "tokens"

type Token struct {
	gorm.Model
	Token       string `form:"token" json:"token" gorm:"index,type:text"`
	SignedTo    string `form:"signed_to" json:"signed_to" gorm:"type:varchar(255)"`
	SignedBy    string `form:"signed_by" json:"signed_by" gorm:"type:varchar(255)"`
	Description string `json:"description" gorm:"type:varchar(255)"`
	Roles       []Role `gorm:"many2many:token_roles;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type TokenList struct {
	ID          int    `json:"id"`
	Token       string `form:"token" json:"token" gorm:"type:text"`
	SignedTo    string `form:"signed_to" json:"signed_to" gorm:"type:varchar(255)"`
	SignedBy    string `form:"signed_by" json:"signed_by" gorm:"type:varchar(255)"`
	Description string `json:"description" gorm:"type:varchar(255)"`
}

func (*TokenList) TableName() string {
	return token_table_name
}

func CreateToken(SignedTo, description string) (enconterError error) {
	var tokenByte string
	if tokenByte, enconterError = token2.SignPermanentToken(SignedTo); enconterError == nil {
		token := &Token{
			SignedTo:    SignedTo,
			SignedBy:    config.CONFIG.AppName,
			Description: description,
			Token:       tokenByte,
		}
		result := DB.Create(token)
		if enconterError = result.Error; enconterError == nil {
			return nil
		}
	}
	return enconterError
}

func GetTokens(title string, offset, limit int, sort string) (map[string]interface{}, error) {
	templates := []*TokenList{}
	response := make(map[string]interface{})
	var count int64
	result := DB.Select([]string{"id", "token", "signed_to", "signed_by", "description"}).
		Limit(limit).Offset(offset).
		Where("deleted_at is ?", nil).
		Where("token like ?", "%"+title+"%").
		Order(token_table_name + ".id " + sort).
		Find(&templates)
	DB.Model(&Token{}).Distinct("id").Count(&count)
	if result.Error != gorm.ErrRecordNotFound {
		response["list"] = templates
		response["total"] = count
		return response, nil
	}
	return nil, result.Error
}

func UpdateTokenWithID(id uint64, signedTo, description string, isUpdate bool) (enconterError error) {
	var token *Token
	if isUpdate && signedTo != "" {
		var tokenByte string
		if tokenByte, enconterError = token2.SignPermanentToken(signedTo); enconterError == nil {
			token = &Token{
				SignedTo:    signedTo,
				Description: description,
				Token:       tokenByte,
			}
		}
	} else {
		token = &Token{
			Description: description,
		}
	}
	if enconterError == nil {
		result := DB.Model(&Token{}).Where("id = ?", id).Updates(token)
		if enconterError = result.Error; enconterError == nil {
			return nil
		}
	}
	return enconterError
}

func DeleteTokenWithID(id uint64) error {
	result := DB.Delete(&Token{}, id)
	if result.Error == nil {
		return nil
	}
	return result.Error
}

func TokenIsDestoryed(tokenStr string) bool {
	var count int64
	result := DB.Model(&Token{}).
		Where("token = ?", tokenStr).
		Where(token_table_name+".deleted_at is NOT ?", nil).
		Unscoped().
		Count(&count)

	if (result.Error != gorm.ErrRecordNotFound || result.Error == nil) && count > 0 {
		return true
	}
	return false
}
