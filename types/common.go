package types

import (
	"reflect"
	"strconv"
	"time"

	"github.com/mitchellh/mapstructure"
)

type Timestamp struct {
	time.Time
}

func (t *Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t.Time).Unix(), 10)), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	i, err := strconv.ParseInt(string(data[:]), 10, 64)
	if err != nil {
		return err
	}
	t = &Timestamp{
		time.Unix(i, 0),
	}
	return nil
}

func Float64ToTimeHookFunc() mapstructure.DecodeHookFunc {
	return func(
		f reflect.Type,
		t reflect.Type,
		data interface{}) (interface{}, error) {
		if f.Kind() != reflect.Float64 {
			return data, nil
		}
		if t != reflect.TypeOf(Timestamp{}) {
			return data, nil
		}

		i := int64(data.(float64))
		return &Timestamp{
			time.Unix(i, 0),
		}, nil
	}
}
