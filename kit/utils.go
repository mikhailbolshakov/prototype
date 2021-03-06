package kit

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"io"
	"time"
)

func NewId() string {
	//return strings.Replace(uuid.NewV4().String(), "-", "", -1)
	return uuid.NewV4().String()
}

func UUID(size int) string {
	u := make([]byte, size)
	io.ReadFull(rand.Reader, u)
	return hex.EncodeToString(u)
}

func Nil() string {
	return uuid.Nil.String()
}

// TODO: remove
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

func MustJson(v interface{}) string {
	s, _ := ToJson(v)
	return s
}

func Json(i interface{}) string {
	r, _ := json.Marshal(i)
	return string(r)
}

func MillisFromTime(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func TimeFromMillis(millis int64) time.Time {
	return time.Unix(0, millis*int64(time.Millisecond))
}

