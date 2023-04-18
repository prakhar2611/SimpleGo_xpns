package Utilities

import (
	renderpkg "github.com/unrolled/render"
)

func GetResponse() *renderpkg.Render {
	response := renderpkg.New(renderpkg.Options{StreamingJSON: true})
	return response
}
