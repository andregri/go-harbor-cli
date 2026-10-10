package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"uuid"

	"log/slog"
)

type Project struct {
	ProjectName  string          `json:"project_name"`
	Public       bool            `json:"public"`
	Metadata     ProjectMetadata `json:"metadata"`
	StorageLimit int             `json:"storage_limit"`
}

type ProjectMetadata struct {
	AutoScan string `json:"auto_scan"`
}

type ErrorsResponse struct {
	Errors []ErrorResponse `json:"errors"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type Client struct {
	BaseUrl  string
	Username string
	Password string
	Logger   *slog.Logger
}

func (c Client) Create(p Project) error {
	client := &http.Client{}

	body, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("Error json.Marshal(): %v", err)
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v2.0/projects", c.BaseUrl), bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("Error http.NewRequest() POST /projects: %v", err)
	}
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	reqID := uuid.New().String()
	req.Header.Set("X-Request-Id", reqID)

	c.Logger.Info("Sending request", "method", req.Method, "url", req.URL.String(), "requestID", reqID)

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error client.Do() POST /projects: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusCreated {
		id := resp.Header.Get("location")
		c.Logger.Info(fmt.Sprintf("Project %q (id: %s) created successfully", p.ProjectName, id), "status", resp.StatusCode, "statusText", resp.Status)
	} else {
		c.Logger.Error("POST /api/v2.0/projects returned error", "status", resp.StatusCode, "statusText", resp.Status)

		errorsResp := ErrorsResponse{}
		decoder := json.NewDecoder(resp.Body)
		if err := decoder.Decode(&errorsResp); err != nil {
			return fmt.Errorf("Error decoding error response: %v", err)
		}
		for _, e := range errorsResp.Errors {
			return fmt.Errorf("%s - %s\n", e.Code, e.Message)
		}
	}

	return nil
}
