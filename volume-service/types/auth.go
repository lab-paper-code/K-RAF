package types

import (
	"encoding/base64"
	"fmt"
)

func GetBase64String(str string) string {
	return base64.StdEncoding.EncodeToString([]byte(str))
}

func GetAuthKey(id, pwd string) string {
	return GetBase64String(fmt.Sprintf("%s:%s", id, pwd))
}
