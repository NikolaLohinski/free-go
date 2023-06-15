package types

import (
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Timestamp struct {
	time.Time
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.Time.Unix(), 10)), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	epoch, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse int: %w", err)
	}

	*t = Timestamp{
		time.Unix(epoch, 0),
	}

	return nil
}

var Float64ToTimeDecodeHook mapstructure.DecodeHookFunc = func(
	field reflect.Type,
	typ reflect.Type,
	data interface{},
) (interface{}, error) {
	if field.Kind() != reflect.Float64 {
		return data, nil
	}

	if typ != reflect.TypeOf(Timestamp{time.Now()}) {
		return data, nil
	}

	epoch, ok := data.(float64)
	if !ok {
		return nil, fmt.Errorf("%v is not of a float", data)
	}

	return &Timestamp{
		time.Unix(int64(epoch), 0),
	}, nil
}
