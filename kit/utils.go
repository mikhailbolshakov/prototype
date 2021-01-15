package kit

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"time"
)

func NewId() string {
	//return strings.Replace(uuid.NewV4().String(), "-", "", -1)
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

func MillisFromTime(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func TimeFromMillis(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}