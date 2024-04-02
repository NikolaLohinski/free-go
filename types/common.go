package types

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

type Timestamp struct {
	time.Time
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(t.Time.Unix(), 10)), nil
}

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	epoch, err := strconv.ParseInt(string(data), 10, 64)
	if err != nil {
		return fmt.Errorf("failed to parse int: %w", err)
	}

	*t = Timestamp{
		time.Unix(epoch, 0).UTC(),
	}

	return nil
}

type Base64Path string

func (c Base64Path) MarshalJSON() ([]byte, error) {
	return json.Marshal(base64.StdEncoding.EncodeToString([]byte(c))) //nolint:wrapcheck
}

func (c *Base64Path) UnmarshalJSON(data []byte) error {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal cd_path: %w", err)
	}

	result, err := base64.StdEncoding.DecodeString(raw)
	if err != nil {
		return fmt.Errorf("failed to decode '%s' from base64: %w", raw, err)
	}

	*c = Base64Path(string(result))

	return nil
}
