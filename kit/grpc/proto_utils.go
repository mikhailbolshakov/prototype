package grpc

import (
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"time"
)

func TimeToPbTS(t *time.Time) *timestamp.Timestamp {

	if t == nil {
		return nil
	} else {
		ts, _ := ptypes.TimestampProto(*t)
		return ts
	}

}

func PbTSToTime(ts *timestamp.Timestamp) *time.Time {

	if ts == nil {
		return nil
	} else {
		t, _ := ptypes.Timestamp(ts)
		return &t
	}
}
