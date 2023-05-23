package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/nikolalohinski/free-go/types"
)

func (c *client) APIVersion() (apiVersion types.APIVersion, err error) {
	response, err := c.httpClient.Get(fmt.Sprintf("%s/api_version", c.base))
	if err != nil {
		err = fmt.Errorf("failed to perform request: %s", err)
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		err = fmt.Errorf("failed to read response body: %s", err)
		return
	}

	if response.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed with status '%d': server returned '%s'", response.StatusCode, string(body))
		return
	}

	if err = json.Unmarshal(body, &apiVersion); err != nil {
		err = fmt.Errorf("failed to unmarshal response body '%s': %s", string(body), err)
		return
	}

	return
}
