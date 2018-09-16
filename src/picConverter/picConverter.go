package picConverter

import (
	"encoding/base64"
)

const base64Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"

var coder = base64.NewEncoding(base64Table)

func Convert2Base64(picBuffer []byte) string {
	return coder.EncodeToString(picBuffer)
}
