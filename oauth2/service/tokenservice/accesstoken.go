package tokenservice

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"github.com/bernishen/lion-go/utils/exception"
	"strings"
)

// SignAccessToken : Sign an access token by client ID, and created a system session.
func SignAccessToken(clientID string) (string, *exception.Exception) {
	var buffer bytes.Buffer
	buffer.WriteString(clientID)
	client := buffer.Bytes()
	if len(client) == 0 {
		return "", exception.NewException(exception.Error, 1001, "The 'clientID' can't be white space.")
	}

	sha := sha256.New()
	sha.Write(client)
	crypher := sha.Sum(nil)

	accesstoken := base64.URLEncoding.EncodeToString(crypher)
	accesstoken = strings.ReplaceAll(accesstoken, "=", "")

	return accesstoken, nil
}
