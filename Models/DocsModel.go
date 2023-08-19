package Models

import (
	"gorm.io/gorm"
)

//database interface
type DocsMeta struct {
	gorm.Model
	ID       string `gorm:"index;unique"`
	MetaData string
	UserId   string
	Title    string `grom:"notnull"`
}

type GetAllDocsResponse struct {
	BaseResponse
	Data []DocsMeta
}
