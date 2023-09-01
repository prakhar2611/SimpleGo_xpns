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
	Folder   string
}

//api request and resposne models

type GetAllDocsResponse struct {
	BaseResponse
	Data []Directory
}

type DocsMetaResponse struct {
	BaseResponse
	Data DocsMeta
}

type GetDocMetaRequest struct {
	Title  string `json:"title"`
	Folder string `json:"folder"`
}

type Directory struct {
	Title    string     `json:"title"`
	Children []Children `json:"children"`
}

type Children struct {
	Title  string `json:"title"`
	Meta   string `json:"meta"`
	Folder string `json:"folder"`
	IsLeaf bool   `json:"isLeaf"`
}
