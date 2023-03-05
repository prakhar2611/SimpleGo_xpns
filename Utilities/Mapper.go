package Utilities

import (
	"SimpleGo_xpns.go/Models"
	"golang.org/x/oauth2"
)

func MapTokenResponse(ID string, token oauth2.Token) *Models.UserToken {
	var t Models.UserToken
	if len(token.AccessToken) > 0 {
		t.ID = ID
		t.AccessToken = token.AccessToken
		t.Expiry = token.Expiry
		t.RefreshToken = token.RefreshToken
		t.TokenType = token.TokenType
		return &t
	}
	return nil
}
