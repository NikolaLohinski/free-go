package types

import (
	"fmt"
	"strconv"
	"time"
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
