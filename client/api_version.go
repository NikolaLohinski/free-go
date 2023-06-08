package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nikolalohinski/free-go/types"
)

func (c *client) APIVersion() (version types.APIVersion, err error) {
	response, err := c.httpClient.Get(fmt.Sprintf("%s/api_version", c.base))
	if err != nil {
		return version, fmt.Errorf("failed to perform request: %w", err)
	}

	defer func() {
		closeError := response.Body.Close()
		if err == nil {
			err = closeError
		} else if closeError != nil {
			err = fmt.Errorf("%w: %w", closeError, err)
		}
	}()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return version, fmt.Errorf("failed to read response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return version, fmt.Errorf(
			"failed with status '%d': server returned '%s'", response.StatusCode, string(body),
		)
	}

	if err = json.Unmarshal(body, &version); err != nil {
		return version, fmt.Errorf("failed to unmarshal response body '%s': %w", string(body), err)
	}

	return version, nil
}
