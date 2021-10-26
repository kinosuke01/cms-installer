package cmsinit

import (
	"encoding/json"
	"errors"
	"fmt"
)

// handle response of init.php
func Handle(statusCode int, bodyBytes []byte) error {
	if statusCode != 200 {
		return fmt.Errorf("status code is %+v", statusCode)
	}

	bodyData := struct {
		Result bool   `json:"result"`
		ErrMsg string `json:"error_message"`
	}{}
	err := json.Unmarshal(bodyBytes, &bodyData)

	if err != nil {
		return err
	}

	if bodyData.ErrMsg != "" {
		return errors.New(bodyData.ErrMsg)
	}

	if !bodyData.Result {
		return errors.New("result is false")
	}

	return nil
}
