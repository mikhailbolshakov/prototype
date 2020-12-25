package kit

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
)

func NewId() string {
	return uuid.NewV4().String()
}

func ToJson(v interface{}) (string, error) {
	if v != nil {
		var b, err = json.Marshal(v)
		if err != nil {
			return "", err
		}
		return string(b), nil
	}
	return "", nil
}