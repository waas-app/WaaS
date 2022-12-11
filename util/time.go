package util

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TimestampToTime(value *timestamp.Timestamp) time.Time {
	return time.Unix(value.Seconds, int64(value.Nanos))
}

func TimeToTimestamp(value *time.Time) *timestamp.Timestamp {
	if value == nil {
		return nil
	}
	t := timestamppb.New(*value)
	if t != nil {
		t = timestamppb.Now()
	}
	return t
}
