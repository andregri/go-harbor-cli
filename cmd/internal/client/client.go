package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	BaseUrl  string
	Username string
	Password string
}

func parseError(resp *http.Response) error {
	fmt.Printf("status code: %d, status message: %s\n", resp.StatusCode, resp.Status)

	errorsResp := ErrorsResponse{}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&errorsResp); err != nil {
		return fmt.Errorf("decoder.Decode(response): %v", err)
	}
	errorsStr := ""
	for _, e := range errorsResp.Errors {
		errorsStr = fmt.Sprintf("%s - %s\n", e.Code, e.Message)
	}
	return fmt.Errorf("%s\n", errorsStr)
}
